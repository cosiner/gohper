package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net"
	netmail "net/mail"
	"net/smtp"
	"strings"

	"github.com/cosiner/gohper/bytes2"
	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gohper/strings2"
	"github.com/cosiner/gohper/unsafe2"
)

const (
	ErrNoTemplate = errors.Err("no template for this type")
)

type mailTemplate struct {
	Subject string
	*template.Template
}

type Mail struct {
	Sender  string
	From    string
	To      []string
	Subject string

	Type string
	Data interface{}

	RawContent string
}

type Mailer struct {
	PrintMail bool
	DoHelo    bool

	addr string
	host string
	auth smtp.Auth

	from   string
	sender string

	Templates  map[string]mailTemplate
	bufferPool bytes2.Pool
	tls        bool
}

func NewMailer(from, sender, username, password, addr string, tls bool) (*Mailer, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	if sender == "" {
		sender = from
	}
	mailer := &Mailer{
		addr:       addr,
		host:       host,
		from:       from,
		sender:     sender,
		bufferPool: bytes2.NewSyncPool(1024, false),
	}
	auth := smtp.PlainAuth("", username, password, strings.Split(addr, ":")[0])
	mailer.auth = auth
	mailer.Templates = make(map[string]mailTemplate)
	mailer.tls = tls
	return mailer, nil
}

func (m *Mailer) AddTemplateFile(typ, filename, subject string) error {
	t, err := template.ParseFiles(filename)
	if err != nil {
		return err
	}

	if typ == "" {
		typ = strings.Split(filename, ".")[0]
	}
	m.Templates[typ] = mailTemplate{
		Subject:  subject,
		Template: t,
	}

	return nil
}

func (m *Mailer) Send(mail *Mail) (err error) {
	tmpl, has := m.Templates[mail.Type]
	if !has && mail.RawContent == "" {
		return ErrNoTemplate
	}

	from := mail.From
	if from == "" {
		from = m.from
	}
	sender := mail.Sender
	if sender == "" {
		sender = m.sender
	}

	buffer := bytes.NewBuffer(m.bufferPool.Get(1024, false))

	buffer.WriteString("To:")
	strings2.WriteStringsToBuffer(buffer, mail.To, ";")

	buffer.WriteString("\r\n")
	nm := netmail.Address{
		Address: from,
		Name:    sender,
	}
	buffer.WriteString("From:" + nm.String() + "\r\n")

	subject := mail.Subject
	if has && subject == "" {
		subject = tmpl.Subject
	}
	buffer.WriteString("Subject:" + subject + "\r\n")
	buffer.WriteString("Content-Type: text/html;charset=UTF-8\r\n\r\n")
	if mail.RawContent != "" {
		buffer.WriteString(mail.RawContent)
	} else {
		err = tmpl.Execute(buffer, mail.Data)
	}

	data := buffer.Bytes()
	if m.PrintMail {
		fmt.Println(unsafe2.String(data))
	}
	if err == nil {
		err = m.send(from, mail.To, data)
	}
	m.bufferPool.Put(data)

	return
}

func (m *Mailer) send(from string, to []string, msg []byte) error {
	if !m.tls {
		return smtp.SendMail(m.addr, m.auth, from, to, msg)
	}

	conn, err := tls.Dial("tcp", m.addr, nil)
	if err != nil {
		return err
	}
	c, err := smtp.NewClient(conn, m.host)
	if err != nil {
		conn.Close()
		return err
	}
	defer c.Close()

	if m.DoHelo {
		err = c.Hello("")
		if err != nil {
			return err
		}
	}
	if m.auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(m.auth); err != nil {
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
