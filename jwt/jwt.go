package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Token struct {
	Raw       string
	Method    string
	Header    map[string]interface{}
	Claims    map[string]interface{}
	Signature byte
	Valid     bool
}

func NewToken(id int, login string, male bool, secret string) string {

	key := []byte(secret)

	h := hmac.New(sha256.New, key)

	header := map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}
	payload := map[string]interface{}{
		"id":    id,
		"login": login,
		"male":  male,
	}
	hdr, err := json.Marshal(header)
	if err != nil {
		fmt.Println("Why isn't it possible?")
	}
	pld, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Why isn't it possible?")
	}
	str := base64.URLEncoding.EncodeToString(hdr) + "." + base64.URLEncoding.EncodeToString(pld)
	h.Write([]byte(str))
	return str + "." + base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(h.Sum(nil))
}
