package crypto

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"strconv"

	"time"
)

var encoding = base64.URLEncoding
var sep = []byte("|")

// SignSecret encode string with given secret key
func SignSecret(secret, value []byte) []byte {
	encVal := make([]byte, encoding.EncodedLen(len(value)))
	encoding.Encode(encVal, value)
	now := make([]byte, 0, 10)
	now = strconv.AppendInt(now, time.Now().UnixNano()<<32>>32, 10)
	hash := hmac.New(sha1.New, secret)
	hash.Write(encVal)
	hash.Write(now)
	sig := hash.Sum(nil)
	buf := bytes.NewBuffer(make([]byte, 0, len(encVal)+len(now)+len(sig)+2))
	buf.Write(encVal)
	buf.Write(sep)
	buf.Write(now)
	buf.Write(sep)
	buf.Write(sig)
	return buf.Bytes()
}

// VerifySecret verify string with given seret key
func VerifySecret(secret, value []byte) (result []byte) {
	var (
		sections = bytes.Split(value, sep)
		encVal   []byte
	)
	if len(sections) == 3 {
		encVal = sections[0]
		now := sections[1]
		sig := sections[2]
		hash := hmac.New(sha1.New, secret)
		hash.Write(encVal)
		hash.Write(now)
		if bytes.Equal(hash.Sum(nil), sig) {
			result = make([]byte, encoding.DecodedLen(len(encVal)))
			if n, e := encoding.Decode(result, encVal); e == nil {
				result = result[:n]
			} else {
				result = nil
			}
		}
	}
	return
}
