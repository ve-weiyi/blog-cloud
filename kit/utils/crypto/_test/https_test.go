package _test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"testing"
	"time"
)

func Test_HTTPS(t *testing.T) {
	// 生成 RSA 密钥对
	privA, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	privB, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	// 创建一个新的证书模板
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// 创建自签名证书
	derBytesA, err := x509.CreateCertificate(rand.Reader, &template, &template, &privA.PublicKey, privA)
	if err != nil {
		log.Fatal(err)
	}
	derBytesB, err := x509.CreateCertificate(rand.Reader, &template, &template, &privB.PublicKey, privB)
	if err != nil {
		log.Fatal(err)
	}

	// PEM 编码证书和密钥
	certOutA := &pem.Block{Type: "CERTIFICATE", Bytes: derBytesA}
	keyOutA := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privA)}
	certOutB := &pem.Block{Type: "CERTIFICATE", Bytes: derBytesB}
	keyOutB := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privB)}

	// 创建 TLS 证书
	certA, err := tls.X509KeyPair(pem.EncodeToMemory(certOutA), pem.EncodeToMemory(keyOutA))
	if err != nil {
		log.Fatal(err)
	}
	certB, err := tls.X509KeyPair(pem.EncodeToMemory(certOutB), pem.EncodeToMemory(keyOutB))
	if err != nil {
		log.Fatal(err)
	}

	// 创建 TLS 配置
	tlsConfigA := &tls.Config{
		Certificates: []tls.Certificate{certA},
	}
	tlsConfigB := &tls.Config{
		Certificates: []tls.Certificate{certB},
	}

	// 打印 TLS 配置以模拟 HTTPS 过程
	log.Println("A's TLS Config:", tlsConfigA)
	log.Println("B's TLS Config:", tlsConfigB)

	// 模拟 A 和 B 互相发送数据
	messageA := "Hello from A"
	messageB := "Hello from B"
	log.Println("A sends:", messageA)
	log.Println("B sends:", messageB)
}
