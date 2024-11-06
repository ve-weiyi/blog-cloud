package _test

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/hex"
	"log"
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/crypto"
)

func hexToBytes(hexStr string) []byte {
	bs, _ := hex.DecodeString(hexStr)
	return bs
}

// 对于使用 AES（Advanced Encryption Standard）算法的 CBC 模式，密钥的长度可以是 128 比特（16 字节）、192 比特（24 字节）或 256 比特（32 字节）。
func Test_AES(t *testing.T) {
	var plaintext []byte // 加密的密钥
	var key []byte       // 加密的密钥
	var iv []byte        // 加密的向量

	plaintext = []byte("791422171@qq.comasdsda")
	key = []byte("1234567.1234567.1234567.1234567.")
	iv = key[:aes.BlockSize] // 初始化向量

	crypt := crypto.AesCBC

	log.Println("iv(hex)", hex.EncodeToString(iv))

	log.Println("秘钥", string(key))
	log.Println("秘钥([]byte)", key)
	log.Println("秘钥(hex)", hex.EncodeToString(key)) // 88403ab936164b48ed66b14fcd864061
	log.Println("秘钥(base64)", base64.StdEncoding.EncodeToString(key))
	log.Println("原文：", string(plaintext))
	log.Println("------------------ CBC模式 --------------------")
	encrypted, err := crypt.AESEncrypt(plaintext, key, iv...)
	log.Println(err)
	log.Println("密文：", encrypted)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
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
	// key = []byte("1234567.1234567.1234567.1234567.")
	copy(iv, key[:aes.BlockSize])

	// iv := slices.Clone(key[:aes.BlockSize])
	// iv := key[:aes.BlockSize] //会导致key被改变
	crypt := crypto.AesCBC

	log.Println("秘钥", key)
	log.Println("秘钥(hex)", hex.EncodeToString(key))
	log.Println("秘钥(base64)", base64.StdEncoding.EncodeToString(key))
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
	data := `yD8cbvnB8b4xjNm/h+BUNr+rztZvnZALS+UXPDEQpEs=`
	key, _ := hex.DecodeString("88403ab936164b48ed66b14fcd864061")

	ie, _ := base64.StdEncoding.DecodeString(data)
	encrypted := ie[aes.BlockSize:]
	iv := ie[:aes.BlockSize] // 初始化向量

	log.Println("秘钥(hex)", hex.EncodeToString(key))
	log.Println("iv(hex)", hex.EncodeToString(iv))

	log.Println("密文：", string(encrypted))
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(append(iv, encrypted...)))
	crypt := crypto.AesCBC
	decrypted, err := crypt.AESDecrypt(encrypted, key, []byte(iv)...)
	log.Println("解密结果：", string(decrypted), err)
}
