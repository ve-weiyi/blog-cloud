package _test

import (
	"bytes"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"testing"
)

func Test_ECDH(t *testing.T) {
	// 选择一个椭圆曲线
	curve := ecdh.P256()

	// 生成通信方A的密钥对
	aliceKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	bobKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	// 通信方A使用自己的私钥和通信方B的公钥生成共享密钥
	bobSecret, err := bobKey.ECDH(aliceKey.PublicKey())
	if err != nil {
		t.Fatal(err)
	}
	aliceSecret, err := aliceKey.ECDH(bobKey.PublicKey())
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(bobSecret, aliceSecret) {
		t.Error("two ECDH computations came out different")
	}
}

func Test_ECDH_ECC(t *testing.T) {
	// 选择一个椭圆曲线
	curve := elliptic.P256()

	// 生成通信方A的密钥对
	privateKeyA, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 生成通信方B的密钥对
	privateKeyB, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 通信方A使用自己的私钥和通信方B的公钥生成共享密钥
	sharedKeyA, _ := privateKeyA.PublicKey.Curve.ScalarMult(privateKeyB.PublicKey.X, privateKeyB.PublicKey.Y, privateKeyA.D.Bytes())

	// 通信方B使用自己的私钥和通信方A的公钥生成共享密钥
	sharedKeyB, _ := privateKeyA.PublicKey.Curve.ScalarMult(privateKeyA.PublicKey.X, privateKeyA.PublicKey.Y, privateKeyB.D.Bytes())

	// 检查两个共享密钥是否相同
	if sharedKeyA.Cmp(sharedKeyB) == 0 {
		fmt.Println("Shared keys match")
	} else {
		fmt.Println("Shared keys do not match")
	}
}
