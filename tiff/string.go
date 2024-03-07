package tiff

import (
	_ "embed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/XC-Zero/zero_common/string/paint"
	"github.com/pkg/errors"
)

type StringOptions struct {
	ImageWidth      uint32
	ImageLength     uint32
	StartIndexX     uint32
	StartIndexY     uint32
	Dpi             uint32
	BackGroundColor color.Color
	FontColor       color.Color
	FontSize        uint32
}

var defaultOptions = StringOptions{
	ImageWidth:      500,
	ImageLength:     300,
	StartIndexX:     5,
	StartIndexY:     50,
	Dpi:             72,
	BackGroundColor: color.White,
	FontColor:       color.Black,
	FontSize:        30,
}

//go:embed  de.ttf
var ded []byte

func StringToRGBA(text string, options ...StringOptions) ([]byte, error) {
	var op = defaultOptions
	if len(options) > 0 {
		op = options[0]
	}
	img := image.NewNRGBA(image.Rect(0, 0, int(op.ImageWidth), int(op.ImageLength)))
	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)
	ctx := paint.NewContext()
	pt := paint.Pt(int(op.StartIndexX), int(op.StartIndexY))
	ctx.SetFontSize(float64(op.FontSize))
	ctx.SetDPI(float64(op.Dpi))
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	font, err := paint.ParseFont(ded)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ctx.SetFont(font)
	ctx.SetSrc(image.Black)

	_, err = ctx.DrawString(text, pt)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	create, err := os.Create(text + ".png")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = png.Encode(create, img)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return img.Pix, nil
}
