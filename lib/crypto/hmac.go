package crypto

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"

	"github.com/cosiner/gohper/lib/types"

	"time"
)

var encoding = base64.URLEncoding
var sep = []byte("|")

// SignSecret encode string with given secret key
func SignSecret(secret []byte, value string) string {
	encVal := make([]byte, encoding.EncodedLen(len(value)))
	encoding.Encode(encVal, types.UnsafeBytes(value))
	now := make([]byte, 0, 10)
	now = strconv.AppendInt(now, time.Now().UnixNano()<<32>>32, 10)
	hash := hmac.New(sha256.New, secret)
	hash.Write(encVal)
	hash.Write(now)
	sig := hash.Sum(nil)
	buf := bytes.NewBuffer(make([]byte, 0, len(encVal)+len(now)+len(sig)+2))
	buf.Write(encVal)
	buf.Write(sep)
	buf.Write(now)
	buf.Write(sep)
	buf.Write(sig)
	return encoding.EncodeToString(buf.Bytes())
}

// VerifySecret verify string with given seret key
func VerifySecret(secret []byte, value string) (result string) {
	if valueBytes, err := encoding.DecodeString(value); err == nil {
		if sections := bytes.Split(valueBytes, sep); len(sections) == 3 {
			encVal := sections[0]
			now := sections[1]
			sig := sections[2]
			hash := hmac.New(sha256.New, secret)
			hash.Write(encVal)
			hash.Write(now)
			if bytes.Equal(hash.Sum(nil), sig) {
				res := make([]byte, encoding.DecodedLen(len(encVal)))
				if n, e := encoding.Decode(res, encVal); e == nil {
					result = string(res[:n])
				}
			}
		}
	}
	return
}
