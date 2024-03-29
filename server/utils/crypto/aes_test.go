package crypto

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"testing"
)

// 对于使用 AES（Advanced Encryption Standard）算法的 CBC 模式，密钥的长度可以是 128 比特（16 字节）、192 比特（24 字节）或 256 比特（32 字节）。
func Test_AES(t *testing.T) {
	var plaintext []byte // 加密的密钥
	var key []byte       // 加密的密钥
	var iv []byte        // 加密的向量

	plaintext = []byte("791422171@qq.comasdsda")
	key = []byte("1234567.1234567.1234567.1234567.")
	copy(iv, key[:aes.BlockSize])

	crypt := AesCBC

	fmt.Println("秘钥", key)
	fmt.Println("秘钥(hex)", hex.EncodeToString(key))
	fmt.Println("iv(hex)", hex.EncodeToString(iv))
	log.Println("原文：", string(plaintext))
	log.Println("------------------ CBC模式 --------------------")
	encrypted, err := crypt.AESEncrypt(plaintext, key, iv...)
	log.Println(err)
	log.Println("密文：", string(encrypted))
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(append(iv, encrypted...)))
	decrypted, err := crypt.AESDecrypt(encrypted, key, iv...)
	log.Println(err)
	log.Println("解密结果：", string(decrypted))
}

func TestAesCFBImpl_AESEncrypt(t *testing.T) {
	var plaintext []byte // 加密的密钥
	var key []byte       // 加密的密钥
	var iv []byte        // 加密的向量

	plaintext = []byte("791422171@qq.comasdsda")
	key, _ = hex.DecodeString("313233343536372e313233343536372e313233343536372e313233343536372e")
	//key = []byte("1234567.1234567.1234567.1234567.")
	copy(iv, key[:aes.BlockSize])

	//iv := slices.Clone(key[:aes.BlockSize])
	//iv := key[:aes.BlockSize] //会导致key被改变
	crypt := AesCBC

	fmt.Println("秘钥", key)
	fmt.Println("秘钥(hex)", hex.EncodeToString(key))
	fmt.Println("iv(hex)", hex.EncodeToString(iv))
	log.Println("原文：", string(plaintext))
	log.Println("------------------ CBC模式 --------------------")
	encrypted, err := crypt.AESEncrypt(plaintext, key, iv...)
	log.Println(err)
	log.Println("密文：", string(encrypted))
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(append(iv, encrypted...)))
	decrypted, err := crypt.AESDecrypt(encrypted, key, iv...)
	log.Println(err)
	log.Println("解密结果：", string(decrypted))
}

func TestAesCFBImpl_AESDecrypt(t *testing.T) {
	ciphertext := `ZjbAUn1Zy2+1hOa62fW4Z7zrDEqaGRptajqapB1hl2neTykHhqPTb3sMATM10h0nQ2jAH5NRbOW8uFb3NFu7wpxV8mXuCHtS+3UH0Ec3/Pw=`
	key, _ := hex.DecodeString("6636C0527D59CB6FB584E6BAD9F5B867BA24CB5006973D97D451ECD27C3C9E30")

	ie, _ := base64.StdEncoding.DecodeString(ciphertext)
	encrypted := ie[aes.BlockSize:]
	iv := key[:aes.BlockSize] // 初始化向量

	log.Println("秘钥(hex)", hex.EncodeToString(key))
	log.Println("iv(hex)", hex.EncodeToString(iv))

	log.Println("密文：", string(encrypted))
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(append(iv, encrypted...)))
	crypt := AesCBC
	decrypted, err := crypt.AESDecrypt(encrypted, key, []byte(iv)...)
	log.Println("解密结果：", string(decrypted), err)
}
