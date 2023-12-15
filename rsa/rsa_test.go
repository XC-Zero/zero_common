package rsa

import (
	"log"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	pk, k, err := GenerateKey()
	if err != nil {
		panic(err)
	}
	var hello = "HELLO?"
	key, err := EncodeByPublicKey(hello, pk)
	if err != nil {
		panic(err)
	}
	log.Println(key)

	text, err := DecodeByPrivateKey(key, k)
	if err != nil {
		panic(err)
	}
	log.Println(text)
	log.Println("?????")

}
