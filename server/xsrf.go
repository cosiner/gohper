package server

import (
	"sync"
	"time"

	"github.com/cosiner/golib/crypto"
)

type (
	XsrfErrorHandler interface {
		HandleXsrfError(req Request, resp Response)
	}

	Xsrf interface {
		Start(interval int) // interval is by seconds
		Stop()
		Set(Response) string
		IsValid(Request) bool
	}

	emptyXsrf struct{}

	XsrfTokenGenerator func() string

	xsrf struct {
		value         string
		valueGen      XsrfTokenGenerator
		valueLifetime int
		*sync.RWMutex
		running bool
		stop    chan bool
	}
)

func (emptyXsrf) Set(Response) string  { return "" }
func (emptyXsrf) IsValid(Request) bool { return true }
func (emptyXsrf) Start(int)            {}
func (emptyXsrf) Stop()                {}

func NewXsrf(valueGen XsrfTokenGenerator, lifetime int) Xsrf {
	if valueGen == nil {
		valueGen = GenXsrfToken
	}
	return &xsrf{
		value:         valueGen(),
		valueLifetime: lifetime,
		RWMutex:       new(sync.RWMutex),
		stop:          make(chan bool, 1),
		running:       true,
	}
}

func (x *xsrf) Stop() {
	if x.running {
		x.running = false
		x.stop <- true
	}
}

func (x *xsrf) Start(interval int) {
	go func() {
		c := time.NewTicker(time.Duration(interval) * time.Second)
		for {
			select {
			case <-c.C:
				x.Lock()
				x.value = x.valueGen()
				x.Unlock()
			case <-x.stop:
				return
			}
		}
	}()
}

func (x *xsrf) Set(resp Response) (value string) {
	if x.running {
		x.RLock()
		value = x.value
		x.RUnlock()
		resp.SetCookieWithExpire(XSRF_NAME, value, x.valueLifetime)
	}
	return
}

func (x *xsrf) IsValid(req Request) (value bool) {
	value = true
	if x.running {
		cookieValue := req.Cookie(XSRF_NAME)
		formValue := req.Param(XSRF_NAME)
		if formValue == "" {
			if formValue = req.Header(HEADER_XSRFTOKEN); formValue == "" {
				formValue = req.Header(HEADER_CSRFTOKEN)
			}
		}
		value = (cookieValue == formValue)
	}
	return
}

func GenXsrfToken() string {
	s, err := crypto.RandAlphanumeric(32)
	if err != nil {
		s = XSRF_ONERROR_TOKEN
	}
	return s
}
