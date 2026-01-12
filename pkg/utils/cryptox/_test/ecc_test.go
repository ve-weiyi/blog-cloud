package _test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/cryptox"
)

func TestECC(t *testing.T) {
	// 选择椭圆曲线参数（这里选择了 P-256 曲线）
	curve := elliptic.P256()

	// 生成密钥对
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		fmt.Println("Failed to generate private key:", err)
		return
	}

	// 显示生成的私钥
	fmt.Println("Private key:")
	fmt.Printf("%x\n", privateKey.D)

	// 获取公钥
	publicKey := &privateKey.PublicKey

	// 显示生成的公钥
	fmt.Println("Public key:")
	fmt.Printf("X: %x\n", publicKey.X)
	fmt.Printf("Y: %x\n", publicKey.Y)

	// 待加密的消息
	message := []byte("hello, world")

	// 加密消息
	ciphertext, err := cryptox.EcdsaEncrypt(publicKey, message)
	if err != nil {
		fmt.Println("Encryption failed:", err)
		return
	}

	// 显示加密后的消息
	fmt.Println("Ciphertext:")
	fmt.Printf("%x\n", ciphertext)

	// 解密消息
	plaintext, err := cryptox.EcdsaDecrypt(privateKey, ciphertext)
	if err != nil {
		fmt.Println("Decryption failed:", err)
		return
	}

	// 显示解密后的消息
	fmt.Println("Plaintext:", string(plaintext))
}
