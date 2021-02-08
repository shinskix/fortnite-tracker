package main

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"image/png"
	"io"
	"log"
	"strings"
)

func drawText(canvas *image.RGBA, text string) {
	lines := strings.Split(text, "\n")
	fgColor := image.Black
	fontFace, err := freetype.ParseFont(gomono.TTF)
	if err != nil {
		log.Fatal(err)
	}
	fontSize := 16.0
	fontDrawer := &font.Drawer{
		Dst: canvas,
		Src: fgColor,
		Face: truetype.NewFace(fontFace, &truetype.Options{
			Size:    fontSize,
			DPI:     72.0,
			Hinting: font.HintingNone,
		}),
	}

	textBounds, _ := fontDrawer.BoundString("|TEXT|")
	textHeight := textBounds.Max.Y - textBounds.Min.Y
	overallTextHeight := textHeight.Mul(fixed.I(len(lines)))
	xPosition := (fixed.I(canvas.Rect.Max.X) - fontDrawer.MeasureString(lines[1])) / 2

	for idx, line := range lines {
		fontDrawer.Dot = fixed.Point26_6{
			X: xPosition,
			Y: fixed.I((canvas.Rect.Max.Y)-overallTextHeight.Ceil())/2 + fixed.I(idx*textHeight.Ceil()),
		}
		fontDrawer.DrawString(line)
	}
}

func CreateImage(out io.Writer, text string) error {
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: image.White}, image.Point{}, draw.Src)
	drawText(img, strings.TrimSpace(text))
	return png.Encode(out, img)
}
