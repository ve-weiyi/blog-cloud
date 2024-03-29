package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

//golang crypt包的AES加密函数的使用 https://www.jianshu.com/p/47e8c137ecd4

type Aes interface {
	AESEncrypt(plaintext []byte, key []byte, iv ...byte) (ciphertext []byte, err error)
	AESDecrypt(ciphertext []byte, key []byte, iv ...byte) (plaintext []byte, err error)
}

var (
	AesECB = aesECBImpl{}
	AesCBC = aesCBCImpl{}
	AesCFB = aesCFBImpl{}
	AesGCM = aesGCMImpl{}
)

/*
1、Electronic Code Book(ECB) 电子密码本模式
ECB模式是最早采用和最简单的模式，相同的明文将永远加密成相同的密文。无初始向量，容易受到密码本重放攻击，一般情况下很少用.
它将加密的数据分成若干组，每组的大小跟加密密钥长度相同，然后每组都用相同的密钥进行加密。
*/
type aesECBImpl struct {
}

func (s *aesECBImpl) AESEncrypt(plaintext []byte, key []byte, iv ...byte) (ciphertext []byte, err error) {
	// AES加密算法的加密块必须是16字节(128bit)，所以不足部分需要填充，常用的填充算法是PKCS7。
	plaintext, err = pkcs7Padding(plaintext, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	// 创建密文接收区
	ciphertext = make([]byte, len(plaintext))

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 分组分块加密
	for bs, be := 0, block.BlockSize(); bs < len(plaintext); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Encrypt(ciphertext[bs:be], plaintext[bs:be])
	}

	return ciphertext, err
}

func (s *aesECBImpl) AESDecrypt(ciphertext []byte, key []byte, iv ...byte) (plaintext []byte, err error) {
	// 创建明文接收区
	plaintext = make([]byte, len(ciphertext))

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 分组分块解密
	for bs, be := 0, block.BlockSize(); bs < len(ciphertext); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Decrypt(plaintext[bs:be], ciphertext[bs:be])
	}

	plaintext, err = pkcs7UnPadding(plaintext)
	return plaintext, err
}

/*
2.Cipher Block Chaining(CBC)  密码分组链接模式
明文被加密前要与前面的密文进行异或运算后再加密，因此只要选择不同的初始向量，相同的密文加密后会形成不同的密文，这是目前应用最广泛的模式。
这种模式是先将明文切分成若干小段，然后每一小段与初始块或者上一段的密文段进行异或运算后，再与密钥进行加密。
*/
type aesCBCImpl struct {
}

func (s *aesCBCImpl) AESEncrypt(plaintext []byte, key []byte, iv ...byte) (ciphertext []byte, err error) {
	// AES加密算法的加密块必须是16字节(128bit)，所以不足部分需要填充，常用的填充算法是PKCS7。
	plaintext, err = pkcs7Padding(plaintext, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	// 创建密文接收区
	ciphertext = make([]byte, len(plaintext))

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//加密向量,取密钥前16位
	if len(iv) == 0 {
		iv = key[:aes.BlockSize]
	}
	if len(iv) != block.BlockSize() {
		return nil, fmt.Errorf("invalid iv '%s' as it's not multiple of ase.blockSize", iv)
	}
	//使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//执行加密
	blockMode.CryptBlocks(ciphertext, plaintext)
	return ciphertext, err
}

func (s *aesCBCImpl) AESDecrypt(ciphertext []byte, key []byte, iv ...byte) (plaintext []byte, err error) {
	// 创建明文接收区
	plaintext = make([]byte, len(ciphertext))

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//加密向量,取密钥前16位
	if len(iv) == 0 {
		iv = key[:aes.BlockSize]
	}
	if len(iv) != block.BlockSize() {
		return nil, fmt.Errorf("invalid iv '%s' as it's not multiple of ase.blockSize", iv)
	}
	//使用cbc
	blockMode := cipher.NewCBCDecrypter(block, iv)
	//执行解密
	blockMode.CryptBlocks(plaintext, ciphertext)
	//去填充
	plaintext, err = pkcs7UnPadding(plaintext)

	return plaintext, err
}

/*
3.Cipher Feedback(CFB)  密码反馈模式
类似于自同步序列密码，分组加密后，按8位分组将密文和明文进行移位异或后得到输出同时反馈回移位寄存器，优点最小可以按字节进行加解密，也可以是n位的，
CFB也是上下文相关的，CFB模式下，明文的一个错误会影响后面的密文(错误扩散)。
*/
type aesCFBImpl struct {
}

func (s *aesCFBImpl) AESEncrypt(plaintext []byte, key []byte, iv ...byte) (ciphertext []byte, err error) {
	// AES加密算法的加密块必须是16字节(128bit)，所以不足部分需要填充，常用的填充算法是PKCS7。
	plaintext, err = pkcs7Padding(plaintext, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	// 创建密文接收区
	ciphertext = make([]byte, len(plaintext))

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//加密向量,取密钥前16位
	if len(iv) == 0 {
		iv = key[:aes.BlockSize]
	}
	if len(iv) != block.BlockSize() {
		return nil, fmt.Errorf("invalid iv '%s' as it's not multiple of ase.blockSize", iv)
	}
	blockMode := cipher.NewCFBEncrypter(block, iv)
	blockMode.XORKeyStream(ciphertext, plaintext)
	return ciphertext, err
}

func (s *aesCFBImpl) AESDecrypt(ciphertext []byte, key []byte, iv ...byte) (plaintext []byte, err error) {
	// 创建明文接收区
	plaintext = make([]byte, len(ciphertext))

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//加密向量,取密钥前16位
	if len(iv) == 0 {
		iv = key[:aes.BlockSize]
	}
	if len(iv) != block.BlockSize() {
		return nil, fmt.Errorf("invalid iv '%s' as it's not multiple of ase.blockSize", iv)
	}
	blockMode := cipher.NewCFBDecrypter(block, iv)
	blockMode.XORKeyStream(plaintext, ciphertext)

	plaintext, err = pkcs7UnPadding(plaintext)
	return plaintext, err
}

/*
4. GCM(Galois/Counter Mode)  基于计数器的模式
GCM模式是一种高级模式，它将加密和完整性校验组合在一起，可以同时实现加密和完整性校验，而且还可以并行计算，因此在速度上比CBC模式更有优势。
*/
type aesGCMImpl struct {
}

/*
key:加密key
plaintext：加密明文 (GCM不需要加密块必须16字节长度，可以是任意长度，其他的都需要16字节对其，所以不足部分都需要补充)
ciphertext:解密返回字节字符串[ 整型以十六进制方式显示]
noncetext: 当前的mac
*/
func (s *aesGCMImpl) AESEncrypt(plaintext []byte, key []byte, iv ...byte) (ciphertext []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := iv
	if len(nonce) == 0 {
		nonce = key[:aesgcm.NonceSize()]
	}
	//nonce := make([]byte, aesgcm.NonceSize())
	//// 由于存在重复的风险，请勿使用给定密钥使用超过2^32个随机值。
	//io.ReadFull(rand.Reader, nonce)

	//out := aesgcm.Seal(nonce, nonce, data, nil) nonce拼接在密文上
	ciphertext = aesgcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, err
}

func (s *aesGCMImpl) AESDecrypt(ciphertext []byte, key []byte, iv ...byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := iv
	if len(nonce) == 0 {
		nonce = key[:aesgcm.NonceSize()]
	}

	plaintext, err = aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, err
}
