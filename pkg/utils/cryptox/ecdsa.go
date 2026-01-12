package cryptox

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

func EcdsaGenerateKey() (string, string, error) {
	// 1. 生成密钥对 (使用Ecdsa曲线)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
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

	return hex.EncodeToString(privateKeyBytes), hex.EncodeToString(publicKeyBytes), nil
}

func EcdsaGenerateSignature(privateKeyStr, hashStr string) (string, error) {
	// 1. 将16进制字符串转换为字节
	privateKeyBytes, err := hex.DecodeString(privateKeyStr)
	if err != nil {
		return "", fmt.Errorf("invalid private key format: %v", err)
	}

	hash, err := hex.DecodeString(hashStr)
	if err != nil {
		return "", fmt.Errorf("invalid hash format: %v", err)
	}

	// 检查输入长度
	if len(privateKeyBytes) != 32 {
		return "", fmt.Errorf("private key must be 32 bytes")
	}
	if len(hash) != 32 {
		return "", fmt.Errorf("hash must be 32 bytes")
	}

	// 解析私钥
	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = elliptic.P256()
	privateKey.D = new(big.Int).SetBytes(privateKeyBytes)
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(privateKeyBytes)

	// 生成签名
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", err
	}

	// 将签名转换为固定长度64字节 (32字节R + 32字节S)
	signature := make([]byte, 64)
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	copy(signature[32-len(rBytes):32], rBytes)
	copy(signature[64-len(sBytes):64], sBytes)

	return hex.EncodeToString(signature), nil
}

func EcdsaVerifySignature(publicKeyStr, hashStr, signatureStr string) (bool, error) {
	// 1. 将16进制字符串转换为字节
	publicKeyBytes, err := hex.DecodeString(publicKeyStr)
	if err != nil {
		return false, fmt.Errorf("invalid public key format: %v", err)
	}

	hash, err := hex.DecodeString(hashStr)
	if err != nil {
		return false, fmt.Errorf("invalid hash format: %v", err)
	}

	signature, err := hex.DecodeString(signatureStr)
	if err != nil {
		return false, fmt.Errorf("invalid signature format: %v", err)
	}

	// 检查输入长度
	if len(publicKeyBytes) != 64 {
		return false, fmt.Errorf("public key must be 64 bytes")
	}

	if len(hash) != 32 {
		return false, fmt.Errorf("hash must be 32 bytes")
	}

	if len(signature) != 64 {
		return false, fmt.Errorf("signature must be 64 bytes")
	}

	// 创建Ecdsa曲线
	curve := elliptic.P256()

	// 从字节中解析公钥X和Y
	x := new(big.Int).SetBytes(publicKeyBytes[:32])
	y := new(big.Int).SetBytes(publicKeyBytes[32:])

	// 检查点是否在曲线上
	if !curve.IsOnCurve(x, y) {
		return false, fmt.Errorf("public key point is not on the curve")
	}

	// 创建公钥对象
	publicKey := &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}

	// 从签名中解析R和S
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])

	// 验证签名
	return ecdsa.Verify(publicKey, hash, r, s), nil
}
