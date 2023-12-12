package token

import (
	"encoding/base64"
	"encoding/json"
	"github.com/XC-Zero/zero_common/aes"
)

// Token 令牌
type Token struct {
	StaffID     int
	StaffName   string
	PrivilegeID []int
	UserEmail   string
}

func (t Token) String() string {
	marshal, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	token, err := aes.Encrypt(marshal, getKey())
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(token)
}

func Parse(tokenStr string) (Token, error) {
	decodeString, err := base64.URLEncoding.DecodeString(tokenStr)
	if err != nil {
		return Token{}, err
	}
	var token Token
	decrypt, err := aes.Decrypt(decodeString, getKey())
	if err != nil {
		return token, err
	}
	err = json.Unmarshal(decrypt, &token)
	if err != nil {
		return token, err
	}
	return token, nil
}

// TODO 换为从配置文件或环境变量中获取真实的秘钥
func getKey() []byte {

	return []byte("0123456789abcdef0123456789abcdef")
}
