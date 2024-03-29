package initest

import (
	"log"
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/server/utils/crypto"
)

func TestBcrypt(t *testing.T) {
	log.Println(crypto.BcryptHash("admin@qq.com"))

	log.Println(crypto.BcryptCheck(
		"admin@qq.com",
		"$2a$10$0KfTcpR186TjQh1yWZ4preGNDO9AqrcW2JB2jI.dx/UPqrTONkvQG"))

	log.Println(crypto.BcryptCheck(
		"admin@qq.com",
		"$2a$10$ZINovpDg.FxFQRj6nhKDLOH55k19RDViybnVVn5EGuKQAcqChRs1e"))
}
