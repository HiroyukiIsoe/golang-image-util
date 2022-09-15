package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	"io/ioutil"

	// _ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/esimov/stackblur-go"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func main() {
	imgFile, err := os.Open("assets/images/go_front.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		println(err.Error())
		return
	}
	
	drawImg := image.NewRGBA(image.Rect(0, 0, 1200, 630))

	sx := 600 - 280
	ex := 600 + 280
	sy := 315 - 280
	ey := 315 + 280

	draw.Draw(drawImg, image.Rect(sx, sy, ex, ey), img, image.Point{0, 0}, draw.Over)

	// 画像ぼかし処理
	blurredImg, err := stackblur.Process(drawImg, uint32(40))
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := saveImgFile(blurredImg, "blurred_out"); err != nil {
		fmt.Println(err)
		return
	}

	if err := drawTextOnImage(drawImg, "Logo Sample.ロゴ サンプル"); err != nil {
		fmt.Println(err)
		return
	}

	if err := saveImgFile(drawImg, "out"); err != nil {
		fmt.Println(err)
		return
	}
}

// TODO:imgの型はimage.RGBAとしているが、image.Image型が望ましい
//      そのためには、imgをdraw.Imageに型変換する必要がある
func drawTextOnImage(img *image.RGBA, text string) error {
	jpFontBin, err := ioutil.ReadFile("assets/fonts/ipaexg.ttf")
	if err != nil {
		return err
	}

	ft, err := truetype.Parse(jpFontBin)
	if err != nil {
		return err
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
		Dst: img,
		Src: image.White,
		Face: face,
		Dot: fixed.Point26_6{},
	}

	fdr.Dot.X = fixed.I(0)
	fdr.Dot.Y = fixed.I(90)

	fdr.DrawString(text)
	return nil
}

func saveImgFile(img image.Image, nm string) error {
	out, err := os.Create("tmp/" + nm + ".jpg")
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
