package imgprocess

import (
	// "os"
	// "log"
	"github.com/nfnt/resize"
	"image"
	"image/color"
)

type Tiler interface {
	MakeTile(img image.RGBA, sourceImgTile image.Image, tileSize int, x int, y int)
}

type AvgColorTiler struct{}
type MCTiler struct{}
type MosaicTiler struct {
	TileImage image.Image
}

//gets the average colour of the tile, and fills in the whole tile with that average color
func (avColorTiler AvgColorTiler) MakeTile(img image.RGBA, sourceImgTile image.Image, tileSize int, x int, y int) {
	r, g, b := GetAvgColor(sourceImgTile)
	c := color.RGBA{r, g, b, 255}
	for i := x; i < x+tileSize; i++ {
		for j := y; j < y+tileSize; j++ {
			img.Set(i, j, c)
		}
	}
}

//based on the explanation by Matt Cutts here: https://www.mattcutts.com/blog/photo-mosaic-effect-with-go/
func (mcTiler MCTiler) MakeTile(img image.RGBA, sourceImgTile image.Image, tileSize int, x int, y int) {
	avgR, avgG, avgB := GetAvgColor(sourceImgTile)
	//my math and casting and everything here is probably soooo bad and wrong, but it works well
	//enough for what I need
	for x1 := sourceImgTile.Bounds().Min.X; x1 < sourceImgTile.Bounds().Max.X; x1++ {
		for y1 := sourceImgTile.Bounds().Min.Y; y1 < sourceImgTile.Bounds().Max.Y; y1++ {
			c := sourceImgTile.At(x1, y1)
			r, g, b, _ := c.RGBA()
			adjustedR := (uint16(avgR) + uint16(r/256)) / 2
			adjustedG := (uint16(avgG) + uint16(g/256)) / 2
			adjustedB := (uint16(avgB) + uint16(b/256)) / 2

			img.Set(x1, y1, color.RGBA{uint8(adjustedR), uint8(adjustedG), uint8(adjustedB), 255})
		}
	}
}

//start of the tiler that will insert pictures to make a real photomosaic
//right now just take the source image, and fill in the tiles with the images
//of itself
func (mosaicTiler MosaicTiler) MakeTile(img image.RGBA, sourceImgTile image.Image, tileSize int, x int, y int) {
	newImage := resize.Resize(uint(tileSize), uint(tileSize), mosaicTiler.TileImage, resize.NearestNeighbor)

	for i := 0; i < tileSize; i++ {
		for j := 0; j < tileSize; j++ {
			r, g, b, _ := newImage.At(i, j).RGBA()
			c := color.RGBA{uint8(r / 256), uint8(g / 256), uint8(b / 256), 255}
			img.Set(x+i, y+j, c)
		}
	}

}
