package grb

import (
	"testing"
)

func TestNewGrb(t *testing.T) {
	_, err := NewGrb("C:\\Users\\XC\\Desktop\\jma_ssm.grb")
	if err != nil {
		panic(err)
	}
}
