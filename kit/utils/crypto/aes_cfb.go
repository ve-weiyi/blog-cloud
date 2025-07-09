package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
)

// CFB 实现 AES 接口
type CFB struct{}

func NewCFB() AES {
	return &CFB{}
}

// AESEncrypt AES-CFB 加密(二进制版)
func (c *CFB) AESEncrypt(plaintext []byte, key []byte, iv ...byte) (ciphertext []byte, err error) {
	// 验证密钥长度
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("invalid key size: must be 16, 24 or 32 bytes")
	}

	// 创建加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 处理 IV
	if len(iv) == 0 {
		iv = key[:aes.BlockSize]
	}
	if len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("invalid iv size: must be %d bytes", aes.BlockSize)
	}

	// 创建密文接收区
	ciphertext = make([]byte, len(plaintext))

	// 使用 CFB 模式加密
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext, plaintext)

	return ciphertext, nil
}

// AESDecrypt AES-CFB 解密(二进制版)
func (c *CFB) AESDecrypt(ciphertext []byte, key []byte, iv ...byte) (plaintext []byte, err error) {
	// 验证密钥长度
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("invalid key size: must be 16, 24 or 32 bytes")
	}

	// 创建加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 处理 IV
	if len(iv) == 0 {
		iv = key[:aes.BlockSize]
	}
	if len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("invalid iv size: must be %d bytes", aes.BlockSize)
	}

	// 创建明文接收区
	plaintext = make([]byte, len(ciphertext))

	// 使用 CFB 模式解密
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

// AESEncryptString AES-CFB 加密(字符串版)
func (c *CFB) AESEncryptHexString(plaintextHex string, keyHex string, ivHex string) (ciphertextHex string, err error) {
	// 解码 hex 字符串
	plaintext, err := hex.DecodeString(plaintextHex)
	if err != nil {
		return "", err
	}

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return "", err
	}

	iv, err := hex.DecodeString(ivHex)
	if err != nil {
		return "", err
	}

	// 加密
	ciphertext, err := c.AESEncrypt(plaintext, key, iv...)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(ciphertext), nil
}

// AESDecryptString AES-CFB 解密(字符串版)
func (c *CFB) AESDecryptHexString(ciphertextHex string, keyHex string, ivHex string) (plaintextHex string, err error) {
	// 解码 hex 字符串
	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return "", err
	}

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return "", err
	}

	iv, err := hex.DecodeString(ivHex)
	if err != nil {
		return "", err
	}

	// 解密
	plaintext, err := c.AESDecrypt(ciphertext, key, iv...)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(plaintext), nil
}
