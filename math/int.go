package math

import (
	"log"
	"strconv"
)

// IntToBinary 十进制转二进制
func IntToBinary(decimal int) string {
	if decimal == 0 {
		return "0"
	}

	binary := ""
	for decimal > 0 {
		remainder := decimal % 2
		binary = strconv.Itoa(remainder) + binary
		decimal = decimal / 2
	}

	return binary
}

// BinaryToHex 二进制转16进制
func BinaryToHex(binary string) string {
	i, err := strconv.ParseInt(binary, 2, 64)
	if err != nil {
		log.Println("无效的二进制数:", binary)
		return ""
	}
	return strconv.FormatInt(i, 16)
}

// HexToBinary 十六进制转二进制
func HexToBinary(hex string) string {
	i, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		log.Println("无效的十六进制数:", hex)
		return ""
	}
	return strconv.FormatInt(i, 2)
}
