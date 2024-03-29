package initest

import (
	"log"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/ve-weiyi/ve-blog-golang/server/global"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/jjwt"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/jsonconv"
)

func TestJwt(t *testing.T) {
	Init()
	tk, err := global.JWT.CreateToken(
		jjwt.TokenExt{
			Uid:      0,
			Username: "121",
		},
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "blog",
		})

	log.Println(tk, err)

	token, err := global.JWT.ParserToken(tk)
	if err != nil {
		return
	}

	log.Println(jsonconv.ObjectToJsonIndent(token))
}
