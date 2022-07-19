package sign

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

func ShaHmac1(source, secret string) string {
	key := []byte(secret)
	hash := hmac.New(sha1.New, key)
	hash.Write([]byte(source))
	signedBytes := hash.Sum(nil)
	signedString := base64.StdEncoding.EncodeToString(signedBytes)
	return signedString
}
