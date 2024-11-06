package jtoken

import (
	"github.com/golang-jwt/jwt/v4"
)

const (
	JwtAudience  = "aud"
	JwtExpire    = "exp"
	JwtId        = "jti"
	JwtIssueAt   = "iat"
	JwtIssuer    = "iss"
	JwtNotBefore = "nbf"
	JwtSubject   = "sub"
	JwtExt       = "ext"
)

type JwtInstance struct {
	SecretKey []byte
}

func NewJwtInstance(SecretKey []byte) *JwtInstance {
	return &JwtInstance{SecretKey}
}

// GenerateJWT 生成JWT
func (that JwtInstance) CreateToken(options ...Option) (string, error) {
	claims := jwt.MapClaims{}

	for _, option := range options {
		option(claims)
	}

	// 创建一个新的JWT token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 设置签名并获取token字符串
	token, err := jwtToken.SignedString(that.SecretKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

// ParseToken 解析JWT
func (that JwtInstance) ParseToken(tokenString string) (*jwt.Token, error) {
	// 解析JWT字符串
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return that.SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
	// 验证token
	//if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	//	return claims, nil
	//}
	//
	//return nil, fmt.Errorf("invalid token")
}
