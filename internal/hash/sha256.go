package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func ComputeSHA256(data []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	res := h.Sum(nil)
	return hex.EncodeToString(res)
}
