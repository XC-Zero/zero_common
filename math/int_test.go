package math

import (
	"log"
	"testing"
)

func TestDecimalToBinary(t *testing.T) {
	log.Println(IntToBinary(555))
	log.Println(HexToBinary("1f5"))
	log.Println(BinaryToHex("000111110101"))
}
