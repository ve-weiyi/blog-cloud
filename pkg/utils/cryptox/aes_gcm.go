package cryptox

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
)

// GCM 实现 AES 接口
type GCM struct{}

func NewGCM() AES {
	return &GCM{}
}

// AESEncrypt AES-GCM 加密(二进制版)
func (g *GCM) AESEncrypt(plaintext []byte, key []byte, iv ...byte) (ciphertext []byte, err error) {
	// 验证密钥长度
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("invalid key size: must be 16, 24 or 32 bytes")
	}

	// 创建加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建 GCM 模式
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 处理 nonce（使用 IV 作为 nonce）
	if len(iv) < aesgcm.NonceSize() {
		return nil, fmt.Errorf("invalid nonce size: must be at least %d bytes", aesgcm.NonceSize())
	}
	nonce := iv[:aesgcm.NonceSize()]

	// 加密数据，将 nonce 和密文拼接在一起
	ciphertext = aesgcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nil
}

// AESDecrypt AES-GCM 解密(二进制版)
func (g *GCM) AESDecrypt(ciphertext []byte, key []byte, iv ...byte) (plaintext []byte, err error) {
	// 验证密钥长度
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("invalid key size: must be 16, 24 or 32 bytes")
	}

	// 创建加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建 GCM 模式
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 处理 nonce（使用 IV 作为 nonce）
	if len(iv) < aesgcm.NonceSize() {
		return nil, fmt.Errorf("invalid nonce size: must be at least %d bytes", aesgcm.NonceSize())
	}
	nonce := iv[:aesgcm.NonceSize()]

	// 解密数据
	plaintext, err = aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// AESEncryptString AES-GCM 加密(字符串版)
func (g *GCM) AESEncryptHexString(plaintextHex string, keyHex string, ivHex string) (ciphertextHex string, err error) {
	// 解码 hex 字符串
	plaintext, err := hex.DecodeString(plaintextHex)
	if err != nil {
		return "", err
	}

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return "", err
	}

	var iv []byte
	if ivHex != "" {
		iv, err = hex.DecodeString(ivHex)
		if err != nil {
			return "", err
		}
	}

	// 加密
	ciphertext, err := g.AESEncrypt(plaintext, key, iv...)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(ciphertext), nil
}

// AESDecryptString AES-GCM 解密(字符串版)
func (g *GCM) AESDecryptHexString(ciphertextHex string, keyHex string, ivHex string) (plaintextHex string, err error) {
	// 解码 hex 字符串
	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return "", err
	}

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return "", err
	}

	var iv []byte
	if ivHex != "" {
		iv, err = hex.DecodeString(ivHex)
		if err != nil {
			return "", err
		}
	}

	// 解密
	plaintext, err := g.AESDecrypt(ciphertext, key, iv...)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(plaintext), nil
}
