package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
)

const (
	RSA_KEY_SIZE = 1024
)

func base64PrivateKey(privateKey *rsa.PrivateKey) string {
	privateBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	return base64.StdEncoding.EncodeToString(privateBytes)
}
func base64PublicKey(publicKey *rsa.PublicKey) (string, error) {
	publicBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(publicBytes), nil
}

func GenerateRsaKeys() (string, string) {
	var public_key, private_key string
	for {
		key, err := rsa.GenerateKey(rand.Reader, RSA_KEY_SIZE)
		if err != nil {
			break
		}
		private_key = base64PrivateKey(key)
		public_key, _ = base64PublicKey(&(key.PublicKey))
		break
	}
	return public_key, private_key
}
