package _test

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/cryptox"
)

func hexToBytes(hexStr string) []byte {
	bs, _ := hex.DecodeString(hexStr)
	return bs
}

// 对于使用 AES（Advanced Encryption Standard）算法的 CBC 模式，密钥的长度可以是 128 比特（16 字节）、192 比特（24 字节）或 256 比特（32 字节）。
func Test_GenerateRandomAESKey(t *testing.T) {
	// 生成AES-256密钥(32字节)
	key, err := cryptox.GenerateRandomAESKey(32)
	if err != nil {
		panic(err)
	}

	fmt.Printf("AES密钥(Hex): %s\n", hex.EncodeToString(key))
	fmt.Printf("AES密钥长度: %d字节\n", len(key))
}

func Test_AES_CBC(t *testing.T) {
	// 生成AES-256密钥(32字节)
	key, err := cryptox.GenerateRandomAESKey(32)
	if err != nil {
		panic(err)
	}

	// 生成随机IV(16字节)
	iv, err := cryptox.GenerateRandomAESKey(16)
	if err != nil {
		panic(err)
	}

	// 明文
	plaintext := []byte("Hello, World! This is a test message for AES encryption.")
	// AES加密
	ciphertext, err := cryptox.NewCBC().AESEncrypt(plaintext, key, iv...)
	if err != nil {
		panic(err)
	}

	// AES解密
	decrypted, err := cryptox.NewCBC().AESDecrypt(ciphertext, key, iv...)
	if err != nil {
		panic(err)
	}

	// 打印结果
	fmt.Println("------------------ Hex --------------------")
	fmt.Printf("密钥(Hex): %s\n", hex.EncodeToString(key))
	fmt.Printf("IV(Hex): %s\n", hex.EncodeToString(iv))
	fmt.Printf("明文(Hex): %s\n", hex.EncodeToString(plaintext))
	fmt.Printf("密文(Hex): %s\n", hex.EncodeToString(ciphertext))

	// base64 编码
	fmt.Println("------------------ base64 --------------------")
	fmt.Printf("密钥(base64): %s\n", base64.StdEncoding.EncodeToString(key))
	fmt.Printf("IV(base64): %s\n", base64.StdEncoding.EncodeToString(iv))
	fmt.Printf("明文(base64): %s\n", base64.StdEncoding.EncodeToString(plaintext))
	fmt.Printf("密文(base64): %s\n", base64.StdEncoding.EncodeToString(ciphertext))

	// 断言解密结果
	if string(decrypted) != string(plaintext) {
		t.Errorf("解密失败: 预期 '%s'，实际 '%s'", plaintext, decrypted)
	} else {
		fmt.Println("解密成功！")
	}
}
