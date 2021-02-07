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
	"log"
	"os"
	"strings"
)

var message = `+--------+------+------------+-------+--------+---------------+
|  MODE  | WINS | WINRATE(%) | KILLS |   KD   |    RATING     |
+--------+------+------------+-------+--------+---------------+
| solo   |    0 |       0.00 |    42 |   1.14 |         1,200 |
| duos   |    9 |       0.70 | 1,150 |   0.84 |         1,093 |
| squads |   63 |       3.90 | 1,747 |   1.14 |         1,190 |
+--------+------+------------+-------+--------+---------------+
|                                      PLAYER | ALIKKLIMENKOV |
+--------+------+------------+-------+--------+---------------+`

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

	rowWithText := lines[1]
	textBounds, _ := fontDrawer.BoundString(rowWithText)
	m := fontDrawer.MeasureString(rowWithText)
	textHeight := textBounds.Max.Y - textBounds.Min.Y
	overallTextHeight := textHeight.Mul(fixed.I(len(lines)))

	for idx, line := range lines {
		xPosition := (fixed.I(canvas.Rect.Max.X) - m) / 2
		yPosition := fixed.I((canvas.Rect.Max.Y)-overallTextHeight.Ceil())/2 + fixed.I(idx*textHeight.Ceil())
		fontDrawer.Dot = fixed.Point26_6{
			X: xPosition,
			Y: yPosition,
		}
		fontDrawer.DrawString(line)
	}
}

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: image.White}, image.Point{}, draw.Src)
	drawText(img, strings.TrimSpace(message))

	file, err := os.Create("test-image.png")
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(file, img)
}
