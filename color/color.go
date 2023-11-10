package color

import (
	"gitlab.tessan.com/data-center/tessan-erp-common/convert"
	"math/rand"
)

func RandColor() string {
	r, g, b := rand.Intn(255), rand.Intn(255), rand.Intn(255)
	return "#" + convert.To16(r) + convert.To16(g) + convert.To16(b)
}
