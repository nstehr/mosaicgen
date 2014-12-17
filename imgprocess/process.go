package imgprocess

import (
	"image"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
)

func GenerateImage(srcImg image.Image, tiler Tiler, tileSize int, outFilename string) {
	width := srcImg.Bounds().Max.X
	height := srcImg.Bounds().Max.Y

	m := image.NewRGBA(image.Rect(0, 0, width, height))
	outFile, err := os.Create(outFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	//checks to see if it has the SubImage method
	img, ok := srcImg.(interface {
		SubImage(r image.Rectangle) image.Image
	})

	if !ok {
		log.Fatal("srcImg does not have SubImage function")
	}

	for x := srcImg.Bounds().Min.X; x < srcImg.Bounds().Max.X; x = x + tileSize {
		for y := srcImg.Bounds().Min.Y; y < srcImg.Bounds().Max.Y; y = y + tileSize {
			subRect := image.Rect(x, y, x+tileSize, y+tileSize)
			sub := img.SubImage(subRect)
			tiler.MakeTile(*m, sub, tileSize, x, y)
		}
	}

	png.Encode(outFile, m)
}

func GetAvgColor(img image.Image) (uint8, uint8, uint8) {
	sumRed := uint32(0)
	sumBlue := uint32(0)
	sumGreen := uint32(0)
	count := uint32(0)
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()

			sumRed += r
			sumBlue += b
			sumGreen += g
			count++

		}
	}

	avgRed := sumRed / uint32(count)
	avgBlue := sumBlue / uint32(count)
	avgGreen := sumGreen / uint32(count)
	//probably bad to do this..hahaha
	return uint8(avgRed / 256), uint8(avgGreen / 256), uint8(avgBlue / 256)

}
