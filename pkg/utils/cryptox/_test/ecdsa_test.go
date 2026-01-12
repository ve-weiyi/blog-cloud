package _test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/cryptox"
)

func Test_(t *testing.T) {
	// 1. 生成密钥对 (使用曲线)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	// 获取私钥字节 (32字节)
	privateKeyBytes := privateKey.D.Bytes()
	// 确保私钥是32字节长度 (可能需要填充前导零)
	if len(privateKeyBytes) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(privateKeyBytes):], privateKeyBytes)
		privateKeyBytes = padded
	}

	// 获取公钥字节 (64字节: X + Y)
	publicKeyBytes := append(
		privateKey.PublicKey.X.Bytes(),
		privateKey.PublicKey.Y.Bytes()...,
	)
	// 确保X和Y都是32字节长度
	if len(publicKeyBytes) != 64 {
		padded := make([]byte, 64)
		copy(padded[32-len(privateKey.PublicKey.X.Bytes()):32], privateKey.PublicKey.X.Bytes())
		copy(padded[64-len(privateKey.PublicKey.Y.Bytes()):], privateKey.PublicKey.Y.Bytes())
		publicKeyBytes = padded
	}

	// 2. 签名数据
	message := []byte("Hello, ECDSA!")
	hash := sha256.Sum256(message)

	// 生成签名
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		panic(err)
	}

	// 将签名转换为固定长度64字节 (32字节R + 32字节S)
	signature := make([]byte, 64)
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	copy(signature[32-len(rBytes):32], rBytes)
	copy(signature[64-len(sBytes):64], sBytes)

	// 打印结果
	printBytes("Private_Key", privateKeyBytes)
	printBytes("Public_Key", publicKeyBytes)
	printBytes("hash", hash[:])
	printBytes("Signature", signature)

	fmt.Println("----")
	// 打印结果
	printHexBytes("Private_Key", privateKeyBytes)
	printHexBytes("Public_Key", publicKeyBytes)
	printHexBytes("hash", hash[:])
	printHexBytes("Signature", signature)

	// 3. 验证签名
	valid := ecdsa.Verify(&privateKey.PublicKey, hash[:], r, s)
	if valid {
		fmt.Println("签名验证成功!")
	} else {
		fmt.Println("签名验证失败!")
	}
}

const (
	privateKeyStr = "1fd6a3dbbc2e201a5e9510bc7fd0436be7156106d9fabdc10cdab909bffd1db1"
	publicKeyStr  = "e5c9b1e188fcd06321df235e3e1ac523db6d52afb1a20fe667cb6a30e41fa9c9d1536e0bf0a7bcc5bc819e10fd2bf19e1db19f4ed3960d0c18a71807d96c6a08"
	hashStr       = "6877498347a58bf169c716d157a503ca85f5f68720d7986a4bd6a9217ad896ca"
	signatureStr  = "f3e5aec55e40f33f58813c4c09c7514830ed2f058a33647d752388a04be51b0586ab53310464e98552ca3aa1eb9c6e986f003c304ef656e099ffb3bff605346e"
)

func Test_GenerateSignature(t *testing.T) {
	//// 1. 将16进制字符串转换为字节
	//privateKeyBytes, err := hex.DecodeString(privateKeyStr)
	//if err != nil {
	//	fmt.Println("公钥格式错误:", err)
	//	return
	//}
	//
	//publicKeyBytes, err := hex.DecodeString(publicKeyStr)
	//if err != nil {
	//	fmt.Println("公钥格式错误:", err)
	//	return
	//}
	//
	//hash, err := hex.DecodeString(hashStr)
	//if err != nil {
	//	fmt.Println("哈希格式错误:", err)
	//	return
	//}
	//
	//signature, err := hex.DecodeString(signatureStr)
	//if err != nil {
	//	fmt.Println("签名格式错误:", err)
	//	return
	//}

	fmt.Println("----", len(privateKeyStr), len(publicKeyStr), len(hashStr), len(signatureStr))

	sign, _ := cryptox.EcdsaGenerateSignature(privateKeyStr, hashStr)
	fmt.Println("----", sign)
	fmt.Println("----", signatureStr)
	fmt.Println("----", sign == signatureStr)
}

func Test_VerifySignature(t *testing.T) {
	// 2. 验证签名
	valid, err := cryptox.EcdsaVerifySignature(publicKeyStr, hashStr, signatureStr)
	if err != nil {
		fmt.Println("验证过程中出错:", err)
		return
	}

	if valid {
		fmt.Println("签名验证成功!")
	} else {
		fmt.Println("签名验证失败!")
	}
}

func Test_printBytes(t *testing.T) {
	privateKeyBytes := []byte{0x2c, 0x13, 0xdb, 0x9c, 0xd5, 0x14, 0x45, 0x20, 0x82, 0x8a, 0xbc, 0xd7, 0x58, 0x2e, 0xfe, 0xaf, 0x13, 0x25, 0x3f, 0x85, 0xa4, 0x12, 0xa2, 0x6a, 0x0d, 0xb1, 0x81, 0x5a, 0x4b, 0x89, 0x28, 0x71}
	publicKeyBytes := []byte{0x79, 0x0e, 0x1e, 0x4c, 0x44, 0x2a, 0xbe, 0x59, 0x4d, 0xcc, 0xac, 0x79, 0xa3, 0x90, 0x84, 0x6e, 0x04, 0xa8, 0x95, 0x8a, 0xf8, 0x38, 0x75, 0x3d, 0x30, 0x1a, 0x4d, 0x94, 0x86, 0xe6, 0xb9, 0xc0, 0xc4, 0x49, 0xf4, 0xf0, 0xfc, 0xd5, 0xea, 0xc3, 0x00, 0x52, 0x13, 0x4c, 0xd5, 0xe6, 0xfb, 0x35, 0xab, 0xdd, 0xf6, 0x0d, 0x1f, 0xb0, 0x15, 0xcd, 0xd9, 0xea, 0x06, 0x97, 0xdc, 0x16, 0xd8, 0xd3}
	hash := []byte{0x68, 0x77, 0x49, 0x83, 0x47, 0xa5, 0x8b, 0xf1, 0x69, 0xc7, 0x16, 0xd1, 0x57, 0xa5, 0x03, 0xca, 0x85, 0xf5, 0xf6, 0x87, 0x20, 0xd7, 0x98, 0x6a, 0x4b, 0xd6, 0xa9, 0x21, 0x7a, 0xd8, 0x96, 0xca}
	signature := []byte{0x29, 0x80, 0x77, 0xa9, 0x07, 0xae, 0x67, 0x13, 0x16, 0xed, 0x51, 0x27, 0xa5, 0xe3, 0x51, 0x5f, 0xe9, 0x84, 0xae, 0xfc, 0xcf, 0xe1, 0xcc, 0x6b, 0xc8, 0x57, 0xc6, 0x5f, 0xe3, 0x6d, 0x9e, 0x51, 0x5f, 0xd0, 0xa8, 0x1f, 0x1d, 0x88, 0x3f, 0xdd, 0xa6, 0x78, 0xb1, 0xfd, 0x6e, 0xd2, 0x92, 0xa6, 0x21, 0xe3, 0xbe, 0x53, 0x8b, 0xcb, 0x2e, 0x04, 0xc1, 0x76, 0xb4, 0x9d, 0x97, 0x34, 0xd2, 0xe4}

	printBytes("privateKeyBytes", privateKeyBytes)
	printBytes("publicKeyBytes", publicKeyBytes)
	printBytes("hash", hash[:])
	printBytes("Signature", signature)
}

// printBytes 打印字节数组为C风格的数组定义
func printBytes(name string, data []byte) {
	fmt.Printf("uint8_t len(%d) %s[] = {", len(data), name)
	for i, b := range data {
		if i > 0 {
			fmt.Print(",")
		}
		fmt.Printf("0x%02x", b)
	}
	fmt.Println("};")
}

func printHexBytes(name string, data []byte) {
	fmt.Println(name, "len(", len(hex.EncodeToString(data)), ") :", hex.EncodeToString(data))
}
