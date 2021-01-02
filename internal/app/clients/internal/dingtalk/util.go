package dingtalk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"
)

// generateSign 加签
func generateSign(secret string) (ts string, sign string) {
	ts = strconv.FormatInt(time.Now().Unix()*1000, 10)
	sign = ts + "\n" + secret

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(sign))

	sign = base64.StdEncoding.EncodeToString(h.Sum(nil))

	return
}
