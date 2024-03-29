package crypto

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"
)

func BenchmarkRSA(b *testing.B) {

	key, _ := rsa.GenerateKey(rand.Reader, RSA_KEY_SIZE)
	pub := &key.PublicKey
	priv := key

	//明文
	plaintext := []byte(Bit512)

	for i := 0; i < b.N; i++ {
		//key-1024bits data-512bits  time-222050 ns/op
		ciphertext, _ := rsa.EncryptOAEP(md5.New(), rand.Reader, pub, plaintext, nil)
		plaintext, _ = rsa.DecryptOAEP(md5.New(), rand.Reader, priv, ciphertext, nil)
	}
}

func TestRSA(t *testing.T) {

	key, _ := rsa.GenerateKey(rand.Reader, RSA_KEY_SIZE)
	pub := &key.PublicKey
	priv := key

	//明文
	plaintext := []byte(Bit512)

	//加密生成密文
	fmt.Printf("%q\n加密:\n", plaintext)
	ciphertext, e := rsa.EncryptOAEP(md5.New(), rand.Reader, pub, plaintext, nil)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Printf("\t%x\n", ciphertext)

	//解密得到明文
	fmt.Printf("解密:\n")
	plaintext, e = rsa.DecryptOAEP(md5.New(), rand.Reader, priv, ciphertext, nil)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Printf("\t%q\n", plaintext)
}
