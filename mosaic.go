package main

import (
	"encoding/json"
	"github.com/nstehr/mosaicgen/db"
	"github.com/nstehr/mosaicgen/imgprocess"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strconv"
	"time"
)

type MosaicMetadata struct {
	Tiles    map[string]imgprocess.Tile
	Height   int
	Width    int
	TileSize int
}

func main() {
	if len(os.Args) <= 3 {
		log.Fatal("not enough args, please pass in source file, keyword, and tile size")
	}
	sourceFile := os.Args[1]
	keyword := os.Args[2]
	tileSize, err := strconv.Atoi(os.Args[3])

	if err != nil {
		log.Fatal("Error parsing tile size")
	}

	in, err := os.Open(sourceFile)
	if err != nil {
		log.Fatal(err)
	}
	srcImg, _, err := image.Decode(in)
	if err != nil {
		log.Fatal(err)
	}

	//dbClient := db.NewMongoClient("localhost")
	//dbClient := db.NewGDatastoreClientFromJSON("", "")
	defer dbClient.CloseConnection()

	tiler := imgprocess.AvgColorTiler{}
	tilerMC := imgprocess.MCTiler{}

	var photos []db.Photo

	dbClient.GetPhotos(keyword, &photos)

	mosaicTiler := imgprocess.NewMosaicTiler(photos)

	imgprocess.GenerateImage(srcImg, tiler, 20, "tiled.png")
	imgprocess.GenerateImage(srcImg, tilerMC, 20, "tiled2.png")
	start := time.Now()
	imgprocess.GenerateImage(srcImg, mosaicTiler, tileSize, "tiled3.png")
	log.Println(time.Since(start))
	//generate some metadata about the mosaic
	metadata := MosaicMetadata{Width: srcImg.Bounds().Max.X, Height: srcImg.Bounds().Max.Y, TileSize: tileSize, Tiles: mosaicTiler.Tiles()}
	metadataJson, _ := json.MarshalIndent(metadata, "", "    ")
	f, err := os.Create("mosaic.json")
	if err != nil {
		log.Fatal("error creating file for mosaic metadata")
	}
	defer f.Close()
	f.Write(metadataJson)
}
