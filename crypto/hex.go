package crypto

import (
	"encoding/hex"
)

// BytesToHexStr transfer binary to hex string
func BytesToHexStr(src []byte) string {
	return hex.EncodeToString(src)
}

// BytesToHex transfer binary to hex bytes
func BytesToHex(src []byte) []byte {
	dst := make([]byte, 2*len(src))
	hex.Encode(dst, src)
	return dst
}

// HexToBytes transfer hex bytes to binary
func HexToBytes(src []byte) []byte {
	dst := make([]byte, len(src)/2)
	hex.Decode(dst, src)
	return dst
}
