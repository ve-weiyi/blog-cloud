package jjwt

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

//1. `Audience`（`aud`）：接收 JWT 的一方。这是一个字符串或字符串数组，表示 JWT 预期的接收者。当应用程序希望指定特定的接收方时，可以使用该字段。
//2. `ExpiresAt`（`exp`）：过期时间。这是一个 Unix 时间戳（以秒为单位），表示 JWT 什么时候会过期。在过期时间之后，JWT 将不再被接受或使用。
//3. `Id`（`jti`）：JWT ID。这是一个用于标识 JWT 的唯一标识符。通常用于防止 JWT 被重复使用。
//4. `IssuedAt`（`iat`）：JWT 的签发时间。这是一个 Unix 时间戳，表示 JWT 什么时候被创建。
//5. `Issuer`（`iss`）：签发者。这是一个字符串，表示 JWT 的签发者。可以用于验证 JWT 的来源是否可信。
//6. `NotBefore`（`nbf`）：在此之前不可用。这是一个 Unix 时间戳，表示 JWT 在什么时间之前不能被接受或使用。
//7. `Subject`（`sub`）：主题。这是一个与 JWT 相关的主题，通常是用户的唯一标识符。它表示 JWT 所涉及的实体。

type TokenClaims struct {
	//Uid       int    `json:"uid"`
	//Username  string `json:"username"`
	//LoginType string `json:"login_type"`
	//Roles    []string `json:"roles"`
	jwt.StandardClaims          //标准荷载,omitempty如果字段为空，则不展示
	Ext                TokenExt `json:"ext,omitempty"` //用户信息,omitempty如果字段为空，则不展示
}

type TokenExt struct {
	Uid       int    `json:"uid"`
	Username  string `json:"username"`
	LoginType string `json:"login_type"`
}

var (
	TokenExpired     = errors.New("token 已过期")
	TokenNotValidYet = errors.New("token 无效")
	TokenMalformed   = errors.New("token 格式错误")
	TokenInvalid     = errors.New("token 不可用")
)

type JwtToken struct {
	SigningKey  []byte        //签名
	Issuer      string        //签发者
	ExpiresTime time.Duration //过期时间
	TokenPrefix string        //token前缀
}

// createToken 生成token
func (j *JwtToken) createToken(claims TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParserToken 解析token
func (j *JwtToken) ParserToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch {
			case ve.Errors&jwt.ValidationErrorMalformed != 0:
				// ValidationErrorMalformed是一个uint常量，表示token不可用
				return nil, TokenMalformed
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				// ValidationErrorExpired表示Token过期
				return nil, TokenExpired
			case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
				// ValidationErrorNotValidYet表示无效token
				return nil, TokenNotValidYet
			default:
				return nil, TokenInvalid
			}
		}
		return nil, err
	}

	if token != nil {
		if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	}

	return nil, TokenInvalid
}

// 根据用户登录信息生成token，
func (j *JwtToken) CreateClaims(userId int, username string, loginType string) (string, error) {
	claims := TokenClaims{
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Add(j.ExpiresTime).Unix(),
			Issuer:    j.Issuer,
		},
	}

	return j.createToken(claims)
}

func (j *JwtToken) CreateToken(ext TokenExt, claims jwt.StandardClaims) (string, error) {
	jwtClaims := TokenClaims{
		Ext:            ext,
		StandardClaims: claims,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString(j.SigningKey)
}

func (j *JwtToken) VerifyToken(token string, uid string) (*TokenClaims, error) {
	if token == "" {
		return nil, errors.New("token is null")
	}
	if uid == "" {
		return nil, errors.New("uid is null")
	}
	//验证token是否 Bearer 开头的
	ok := strings.HasPrefix(token, j.TokenPrefix)
	if !ok {
		return nil, errors.New("token must be has prefix :" + j.TokenPrefix)
	}

	token = strings.TrimPrefix(token, j.TokenPrefix)
	// 解析token
	claims, err := j.ParserToken(token)
	if err != nil {
		return nil, err
	}

	if uid != strconv.Itoa(claims.Ext.Uid) {
		return nil, errors.New("uid is not equal")
	}
	return claims, nil
}
