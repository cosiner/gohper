package tls2

import (
	"crypto/x509"
	"io/ioutil"

	"github.com/cosiner/gohper/errors"
)

const ErrBadPEMFile = errors.Err("pem file can't be parsed")

// CAPool create a ca pool use pem files
func CAPool(pems ...string) (p *x509.CertPool, err error) {
	var data []byte
	for i := 0; i < len(pems) && err == nil; i++ {
		if data, err = ioutil.ReadFile(pems[i]); err == nil {
			if p == nil {
				p = x509.NewCertPool()
			}

			if !p.AppendCertsFromPEM(data) {
				err = ErrBadPEMFile
			}
		}
	}

	return
}
