package aes

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef") // 32字节的AES-256密钥

	plaintext := []byte("你好 mldddd!")

	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return
	}

	fmt.Println("Ciphertext:", base64.StdEncoding.EncodeToString(ciphertext))

	decryptedText, err := Decrypt(ciphertext, key)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}

	fmt.Println("Decrypted text:", string(decryptedText))
}
