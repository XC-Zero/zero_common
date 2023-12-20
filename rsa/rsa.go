package rsa

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/pkg/errors"
)

func GenerateKey() (public, private string, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	k := bytes.NewBufferString("")

	err = pem.Encode(k, privateKeyBlock)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	pk := bytes.NewBufferString("")
	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	}
	err = pem.Encode(pk, publicKeyBlock)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	public = string(pk.Bytes())
	private = string(k.Bytes())
	return
}

func EncodeByPublicKey(plainText, publicKey string, needBase64 bool) (cipherText string, err error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		err = errors.New("无法解析公钥文件" + block.Type)
		return
	}
	pk, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	cipher, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pk, []byte(plainText), nil)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if needBase64 {
		return base64.RawStdEncoding.EncodeToString(cipher), nil
	}
	cipherText = string(cipher)
	return
}

func DecodeByPrivateKey(cipherText, privateKey string, needBase64 bool) (plainText string, err error) {
	if needBase64 {
		txt, err := base64.RawStdEncoding.DecodeString(cipherText)
		if err != nil {
			return "", err
		}
		cipherText = string(txt)

	}
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		err = errors.Errorf("无法解析私钥文件")
		return
	}
	pk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	text, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, pk, []byte(cipherText), nil)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	plainText = string(text)
	return
}
