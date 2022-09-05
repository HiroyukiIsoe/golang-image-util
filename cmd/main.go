package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	"image/jpeg"

	// _ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/esimov/stackblur-go"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"
)

func main() {
	imgFile, err := os.Open("assets/go_front.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer imgFile.Close()

	img, imgFmt, err := image.Decode(imgFile)
	if err != nil {
		println(err.Error())
		return
	}
	
	fmt.Println("Image Format is ", imgFmt)

	drawImg := image.NewRGBA(image.Rect(0, 0, 1200, 630))

	sx := 600 - 280
	ex := 600 + 280
	sy := 315 - 280
	ey := 315 + 280

	fmt.Println(image.Rect(sx, sy, ex, ey))
	draw.Draw(drawImg, image.Rect(sx, sy, ex, ey), img, image.Point{0, 0}, draw.Over)

	// ぼかし画像
	blurredImg, err := stackblur.Process(drawImg, uint32(40))
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := saveImgFile(blurredImg, "blurred_out"); err != nil {
		fmt.Println(err)
		return
	}

	ft, err := truetype.Parse(gobold.TTF)
	if err != nil {
		fmt.Println(err)
		return
	}

	ftOpt := truetype.Options{
		Size: 90,
		DPI: 0,
		Hinting: 0,
		GlyphCacheEntries: 0,
		SubPixelsX: 0,
		SubPixelsY: 0,
	}

	face := truetype.NewFace(ft, &ftOpt)

	fdr := &font.Drawer{
		Dst: drawImg,
		Src: image.White,
		Face: face,
		Dot: fixed.Point26_6{},
	}

	fdr.Dot.X = fixed.I(0)
	fdr.Dot.Y = fixed.I(90)

	fdr.DrawString("2022 Logo Sample😁")

	if err := saveImgFile(drawImg, "blurred_out"); err != nil {
		fmt.Println(err)
		return
	}
}

func saveImgFile(img image.Image, nm string) error {
	out, err := os.Create(nm + ".jpg")
	if err != nil {
		return err
	}
	defer out.Close()

	if err := jpeg.Encode(out, img, nil); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
