package jtoken

import (
	"github.com/golang-jwt/jwt/v4"
)

type Option func(jwt.MapClaims)

func WithAudience(aud string) Option {
	return func(claims jwt.MapClaims) {
		claims["aud"] = aud
	}
}

func WithExpiresAt(exp int64) Option {
	return func(claims jwt.MapClaims) {
		claims["exp"] = exp
	}
}

func WithId(jti string) Option {
	return func(claims jwt.MapClaims) {
		claims["jti"] = jti
	}
}

func WithIssuedAt(iat int64) Option {
	return func(claims jwt.MapClaims) {
		claims["iat"] = iat
	}
}

func WithIssuer(iss string) Option {
	return func(claims jwt.MapClaims) {
		claims["iss"] = iss
	}
}

func WithNotBefore(nbf int64) Option {
	return func(claims jwt.MapClaims) {
		claims["nbf"] = nbf
	}
}

func WithSubject(sub string) Option {
	return func(claims jwt.MapClaims) {
		claims["sub"] = sub
	}
}

func WithClaimExt(key string, value interface{}) Option {
	return func(claims jwt.MapClaims) {
		claims[key] = value
	}
}
