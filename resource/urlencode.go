package resource

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"strings"

	bytes2 "github.com/cosiner/gohper/lib/bytes"
	io2 "github.com/cosiner/gohper/lib/io"
	url2 "github.com/cosiner/gohper/lib/net/url"
	"github.com/cosiner/gohper/lib/reflect"
)

const TAG_URLENCODE = "encode"

type URLEncode struct {
	pool bytes2.Pool
}

func NewURLEncode(pool bytes2.Pool) URLEncode {
	if pool == nil {
		pool = bytes2.NewSyncPool(256, false) // don't allow small buffer
	}
	pool.Init()
	return URLEncode{pool}
}

func (u URLEncode) Marshal(v interface{}) ([]byte, error) {
	switch val := v.(type) {
	case []byte:
		return val, nil
	case string:
		return []byte(escape(val)), nil
	case map[string]string:
		buffer := bytes.NewBuffer(u.pool.Get(0, false))
		buf, used := url2.EscapeEncode(val, buffer)
		if !used {
			u.pool.Put(buffer.Bytes())
		}
		return buf, nil
	case map[string][]string:
		return []byte(url.Values(val).Encode()), nil
	default:
		mp := make(map[string]string)
		reflect.MarshalStruct(val, mp, TAG_URLENCODE)

		buffer := bytes.NewBuffer(u.pool.Get(0, false))
		buf, used := url2.EscapeEncode(mp, buffer)
		if !used {
			u.pool.Put(buffer.Bytes())
		}
		return buf, nil
	}
}

func (URLEncode) Unmarshal(data []byte, v interface{}) error {
	s, err := url.QueryUnescape(io2.String(data))
	if err != nil {
		return err
	}

	vals, err := url.ParseQuery(s)
	if err != nil {
		return err
	}

	switch v := v.(type) {
	case *map[string][]string:
		if *v != nil {
			*v = vals
		} else {
			panic("can't use nil address to accept values")
		}
	case map[string][]string:
		for k := range vals {
			v[k] = vals[k]
		}
	case map[string]string:
		for k := range vals {
			v[k] = strings.Join(vals[k], ",")
		}
	default:
		mp := make(map[string]string)
		for k := range vals {
			mp[k] = strings.Join(vals[k], ",")
		}
		reflect.UnmarshalStruct(v, mp, TAG_URLENCODE)
	}
	return nil
}

func (u URLEncode) Pool(buf []byte) {
	u.pool.Put(buf)
}

func (u URLEncode) Send(w io.Writer, key string, v interface{}) error {
	if key != "" {
		ew := io2.NewErrorWriter(w)
		ew.WriteString(escape(key))
		ew.WriteString(escape("="))
		switch s := v.(type) {
		case string:
			ew.WriteString(escape(s))
		case []byte:
			ew.WriteString(escape(io2.String(s)))
		default:
			ew.WriteString(escape(fmt.Sprint(v)))
		}
		return ew.Error
	}

	buf, err := u.Marshal(v)
	if err == nil {
		_, err = w.Write(buf)
	}
	u.pool.Put(buf)
	return err
}

func (u URLEncode) Receive(r io.Reader, v interface{}) error {
	buf := bytes.NewBuffer(u.pool.Get(0, false))
	defer u.pool.Put(buf.Bytes())

	n, err := buf.ReadFrom(r)
	if n > 0 && err == nil {
		err = u.Unmarshal(buf.Bytes(), v)
	}
	return err
}

var escape = url.QueryEscape
