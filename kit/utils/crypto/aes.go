package crypto

import (
	"crypto/rand"
	"fmt"
)

//模式对比表
//模式	是否需要填充	是否流密码	是否提供认证	并行化	典型用途
//ECB		是			否			否		是		不推荐使用
//CBC		是			否			否		否		文件加密
//CFB		否			是			否		否		流数据加密
//OFB		否			是			否		否		卫星通信
//CTR		否			是			否		是		高性能加密
//GCM		否			是			是		是		TLS, IPsec
//CCM		是			部分			是		否		无线通信

type AES interface {
	AESEncrypt(plaintext []byte, key []byte, iv ...byte) (ciphertext []byte, err error)
	AESDecrypt(ciphertext []byte, key []byte, iv ...byte) (plaintext []byte, err error)
	AESEncryptHexString(plaintextHex string, keyHex string, ivHex string) (ciphertextHex string, err error)
	AESDecryptHexString(ciphertextHex string, keyHex string, ivHex string) (plaintextHex string, err error)
}

// GenerateRandomAESKey 生成随机AES密钥
// keySize: 16(128-bit), 24(192-bit) 或 32(256-bit)
func GenerateRandomAESKey(keySize int) ([]byte, error) {
	if keySize != 16 && keySize != 24 && keySize != 32 {
		return nil, fmt.Errorf("invalid key size: must be 16, 24 or 32 bytes")
	}

	key := make([]byte, keySize)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate random key: %v", err)
	}

	return key, nil
}
