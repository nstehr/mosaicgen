package main

import (
	"github.com/nstehr/mosaicgen/db"
	"github.com/nstehr/mosaicgen/imgprocess"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"time"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("not enough args, please pass in source file")
	}
	sourceFile := os.Args[1]
	in, err := os.Open(sourceFile)
	if err != nil {
		log.Fatal(err)
	}
	srcImg, _, err := image.Decode(in)
	if err != nil {
		log.Fatal(err)
	}

	dbClient := db.NewMongoClient("localhost")
	defer dbClient.CloseConnection()

	tiler := imgprocess.AvgColorTiler{}
	tilerMC := imgprocess.MCTiler{}
	mosaicTiler := imgprocess.MosaicTiler{TileImage: srcImg, Keyword: "snowman", DB: dbClient}

	imgprocess.GenerateImage(srcImg, tiler, 20, "tiled.png")
	imgprocess.GenerateImage(srcImg, tilerMC, 20, "tiled2.png")
	start := time.Now()
	imgprocess.GenerateImage(srcImg, mosaicTiler, 50, "tiled3.png")
	log.Println(time.Since(start))
}
