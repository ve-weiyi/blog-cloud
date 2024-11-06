package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"
)

// 使用公钥加密消息
func EcdsaEncrypt(publicKey *ecdsa.PublicKey, message []byte) ([]byte, error) {
	// 生成临时密钥对
	tempPrivateKey, err := ecdsa.GenerateKey(publicKey.Curve, rand.Reader)
	if err != nil {
		return nil, err
	}

	// 使用公钥加密消息
	ciphertextX, ciphertextY := publicKey.Curve.ScalarMult(publicKey.X, publicKey.Y, tempPrivateKey.D.Bytes())
	ciphertext := append(ciphertextX.Bytes(), ciphertextY.Bytes()...)

	// 添加加密后的消息
	ciphertext = append(ciphertext, message...)

	return ciphertext, nil
}

// 使用私钥解密消息
func EcdsaDecrypt(privateKey *ecdsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	curve := privateKey.Curve
	keySize := curve.Params().BitSize / 8

	// 解密消息
	ciphertextX, ciphertextY := new(big.Int), new(big.Int)
	ciphertextX.SetBytes(ciphertext[:keySize])
	ciphertextY.SetBytes(ciphertext[keySize : 2*keySize])

	plaintext := make([]byte, len(ciphertext)-2*keySize)
	copy(plaintext, ciphertext[2*keySize:])

	return plaintext, nil
}
