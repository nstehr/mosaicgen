package imgprocess

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/nfnt/resize"
	"github.com/nstehr/mosaicgen/db"
	"github.com/nstehr/mosaicgen/kdtree"
	"image"
	"image/color"
	"image/draw"
	"log"
	"net/http"
	"strconv"
)

type Tiler interface {
	MakeTile(img *image.RGBA, sourceImgTile image.Image, tileSize int, x int, y int)
}

type Tile struct {
	Photo    db.Photo
	Position image.Point
}

type AvgColorTiler struct{}
type MCTiler struct{}
type MosaicTiler struct {
	tiles  map[string]Tile
	photos []db.Photo //not sure if this needs to be a pointer, since it is a slice?
	tree   *kdtree.Node
}

//gets the average colour of the tile, and fills in the whole tile with that average color
func (avColorTiler AvgColorTiler) MakeTile(img *image.RGBA, sourceImgTile image.Image, tileSize int, x int, y int) {
	r, g, b := GetAvgColor(sourceImgTile)
	c := color.RGBA{r, g, b, 255}
	draw.Draw(img, image.Rect(x, y, x+tileSize, y+tileSize), &image.Uniform{c}, image.ZP, draw.Src)
}

//based on the explanation by Matt Cutts here: https://www.mattcutts.com/blog/photo-mosaic-effect-with-go/
func (mcTiler MCTiler) MakeTile(img *image.RGBA, sourceImgTile image.Image, tileSize int, x int, y int) {
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

func NewMosaicTiler(sourcePhotos []db.Photo) *MosaicTiler {
	tiles := make(map[string]Tile)
	return &MosaicTiler{tiles: tiles, photos: sourcePhotos, tree: kdtree.NewTree(sourcePhotos)}
}

func (mosaicTiler *MosaicTiler) MakeTile(img *image.RGBA, sourceImgTile image.Image, tileSize int, x int, y int) {

	aR, aG, aB := GetAvgColor(sourceImgTile)
	cf := colorful.Color{float64(aR) / 255.0, float64(aG) / 255.0, float64(aB) / 255.0}
	matchedPhoto, tileImage, err := mosaicTiler.findClosestImage(cf)
	if err != nil {
		log.Fatal("Error locating matching image")
	}
	normalizedX := x / tileSize
	normalizedY := y / tileSize
	//this key will change when I implement the functionality to use an image only once
	//the key will become the URL
	key := strconv.Itoa(normalizedX) + "," + strconv.Itoa(normalizedY)
	mosaicTiler.tiles[key] = Tile{Photo: matchedPhoto, Position: image.Point{X: normalizedX, Y: normalizedY}}
	newImage := resize.Resize(uint(tileSize), uint(tileSize), tileImage, resize.NearestNeighbor)
	draw.Draw(img, image.Rect(x, y, x+tileSize, y+tileSize), newImage, image.ZP, draw.Src)

}

func (mosaicTiler *MosaicTiler) findClosestImage(sourceAvgColor colorful.Color) (db.Photo, image.Image, error) {

	var photoUrl string
	//matchedPhoto := mosaicTiler.tree.NearestNeighbour(sourceAvgColor)
	matchedPhoto := linearSearch(mosaicTiler.photos, sourceAvgColor)
	if matchedPhoto.ThumbUrl != "" {
				photoUrl = matchedPhoto.ThumbUrl
			} else {
				photoUrl = matchedPhoto.Url
			}
	
	resp, err := http.Get(photoUrl)
	if err != nil {
		log.Printf("error retrieving picture from %s: %s\n", photoUrl, err)
		return matchedPhoto, nil, err
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Printf("error decoding picture from %s: %s\n", photoUrl, err)
		return matchedPhoto, nil, err
	}
	return matchedPhoto, img, nil
}

func (mosaicTiler *MosaicTiler) Tiles() map[string]Tile {
	//this will get more implementation when I implement the functionality
	//to use an image only once
	return mosaicTiler.tiles
}

func linearSearch(photos []db.Photo, sourceAvgColor colorful.Color) (db.Photo) {
	var matchedPhoto db.Photo
	minDistance := 1.0
	for _, photo := range photos {
		c := colorful.Color{float64(photo.AvgColor.R) / 255.0, float64(photo.AvgColor.G) / 255.0, float64(photo.AvgColor.B) / 255.0}
		distance := sourceAvgColor.DistanceCIE94(c)
		if distance < minDistance {
			minDistance = distance
			matchedPhoto = photo
		
		}
	}
	return matchedPhoto
}
