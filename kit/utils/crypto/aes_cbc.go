package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
)

// CBC 实现 AES 接口
type CBC struct{}

func NewCBC() AES {
	return &CBC{}
}

// AESEncrypt AES-CBC 加密(二进制版)
func (c *CBC) AESEncrypt(plaintext []byte, key []byte, iv ...byte) (ciphertext []byte, err error) {
	//加密向量,取密钥前16位
	if len(iv) == 0 {
		iv = key[:aes.BlockSize]
	}
	if len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("invalid iv '%s' as it's not multiple of ase.blockSize", iv)
	}

	// 验证密钥长度
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("invalid key size: must be 16, 24 or 32 bytes")
	}

	// 创建加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// AES加密算法的加密块必须是16字节(128bit)，所以不足部分需要填充，常用的填充算法是PKCS7。
	plaintext, err = pkcs7Padding(plaintext, aes.BlockSize)
	if err != nil {
		return nil, err
	}

	ciphertext = make([]byte, len(plaintext))
	//使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(ciphertext, plaintext)
	return ciphertext, err
}

// AESDecrypt AES-CBC 解密(二进制版)
func (c *CBC) AESDecrypt(ciphertext []byte, key []byte, iv ...byte) (plaintext []byte, err error) {
	//加密向量,取密钥前16位
	if len(iv) == 0 {
		iv = key[:aes.BlockSize]
	}
	if len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("invalid iv '%s' as it's not multiple of ase.blockSize", iv)
	}

	// 验证密钥长度
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("invalid key size: must be 16, 24 or 32 bytes")
	}

	// 检查最小长度
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	// 创建解密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintext = make([]byte, len(ciphertext))
	// 解密
	blockMode := cipher.NewCBCDecrypter(block, iv)
	blockMode.CryptBlocks(plaintext, ciphertext)
	// PKCS7 去除填充
	plaintext, err = pkcs7UnPadding(plaintext)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// AESEncryptString AES-CBC 加密(字符串版)
func (c *CBC) AESEncryptHexString(plaintextHex string, keyHex string, ivHex string) (ciphertextHex string, err error) {
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

// AESDecryptString AES-CBC 解密(字符串版)
func (c *CBC) AESDecryptHexString(ciphertextHex string, keyHex string, ivHex string) (plaintextHex string, err error) {
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
