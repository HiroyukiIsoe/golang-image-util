package main

import (
	"flag"
	"fmt"
	"image"
	"image-util/internal/s3"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"os"
	"time"

	"github.com/esimov/stackblur-go"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func main() {
	var originalFilePath string
	var drawText string
	var isBlur bool
	flag.StringVar(&originalFilePath, "file", "", "originalFilePath")
	flag.StringVar(&drawText, "text", "", "drawText")
	flag.BoolVar(&isBlur, "blur", false, "isBlur")
	flag.Parse()

	convertToOgp(originalFilePath, drawText, isBlur)
}

func convertToOgp(inputFile string, drawText string, isBlur bool) {
	imgFile, err := os.Open(inputFile)
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
	
	var drawImg draw.Image = image.NewRGBA(image.Rect(0, 0, 1200, 630))
	position := calcImagePosition(img.Bounds())

	draw.Draw(drawImg, position, img, image.Point{0, 0}, draw.Over)

	if isBlur {
			// 画像ぼかし処理
			drawImg, err = stackblur.Process(drawImg, uint32(40))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if err := drawTextOnImage(drawImg, drawText); err != nil {
		fmt.Println(err)
		return
	}

	now := time.Now()

	name := "out" + now.Format("_20060102150405")

	if err := saveImgFile(drawImg, name); err != nil {
		fmt.Println(err)
		return
	}

	out, err := os.Open("tmp/" + name + ".jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer out.Close()
	s3.Upload(name + ".jpg", out)
}

func calcImagePosition(baseRect image.Rectangle) image.Rectangle {
	dx := baseRect.Dx()
	dy := baseRect.Dy()
	sx := 600 - (dx / 2)
	ex := 600 + (dx / 2)
	sy := 315 - (dy / 2)
	ey := 315 + (dy / 2)

	return image.Rect(sx, sy, ex, ey)
}

func drawTextOnImage(img draw.Image, text string) error {
	jpFontBin, err := os.ReadFile("assets/fonts/ipaexg.ttf")
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
	path := "tmp/" + nm + ".jpg"
	// backUpCopyFile(path)
	out, err := os.Create(path)
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
