package common

import (
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"time"
)

var jwtKey = []byte("e8236bce52ef2b61")

// CreateToken 创建Token
func CreateToken(payload map[string]interface{}) (string, error) {
	token := jwt.New()
	token.Set(jwt.IssuerKey, `qiaowiwi`)
	token.Set(jwt.AudienceKey, `User`)
	token.Set(jwt.IssuedAtKey, time.Now().Unix())
	token.Set(jwt.ExpirationKey, time.Now().AddDate(0,0,30).Unix())

	for s := range payload {
		token.Set(s, payload[s])
	}

	tokenString, err := jwt.Sign(token, jwa.HS512, jwtKey)
	return string(tokenString), err
}

// VerifyToken 验证Token
func VerifyToken(tokenString string) bool {
	token, err := ParseToken(tokenString)
	if err != nil {
		return false
	}
	if token.Issuer() != "qiaowiwi" {
		return false
	}
	if len(token.Audience()) > 0 && token.Audience()[0] != "User" {
		return false
	}
	if token.Expiration().Before(time.Now()) {
		return false
	}
	return true
}

// ParseToken 解密token拿取中间的数据
func ParseToken(tokenString string) (jwt.Token, error) {
	return jwt.ParseString(tokenString, jwt.WithValidate(true), jwt.WithVerify(jwa.HS512, jwtKey))
}