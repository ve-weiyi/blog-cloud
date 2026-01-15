package _test

import (
	"fmt"
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/cryptox"
)

func TestBcrypt(t *testing.T) {
	pwd := "123456"
	fmt.Println(cryptox.BcryptHash(pwd))
}
