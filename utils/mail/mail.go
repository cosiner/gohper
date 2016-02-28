package mail

import (
	"bytes"
	"fmt"
	"html/template"
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
	From    string
	To      []string
	Subject string

	Type string
	Data interface{}

	RawContent string
}

type Mailer struct {
	PrintMail bool

	addr     string
	auth     smtp.Auth
	username string

	Templates  map[string]mailTemplate
	bufferPool bytes2.Pool
}

func NewMailer(username, password, addr string) *Mailer {
	mailer := &Mailer{
		addr:       addr,
		username:   username,
		bufferPool: bytes2.NewSyncPool(1024, false),
	}
	auth := smtp.PlainAuth("", username, password, strings.Split(addr, ":")[0])
	mailer.auth = auth
	mailer.Templates = make(map[string]mailTemplate)

	return mailer
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
		from = m.username
	}

	buffer := bytes.NewBuffer(m.bufferPool.Get(1024, false))

	buffer.WriteString("To:")
	strings2.WriteStringsToBuffer(buffer, mail.To, ";")

	buffer.WriteString("\r\n")
	buffer.WriteString("From:" + from + "\r\n")

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
		err = smtp.SendMail(m.addr, m.auth, from, mail.To, data)
	}
	m.bufferPool.Put(data)

	return
}
