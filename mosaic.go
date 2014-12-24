package main

import (
	"github.com/nstehr/mosaicgen/imgprocess"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
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
	tiler := imgprocess.AvgColorTiler{}
	tilerMC := imgprocess.MCTiler{}
	mosaicTiler := imgprocess.MosaicTiler{TileImage: srcImg}

	imgprocess.GenerateImage(srcImg, tiler, 20, "tiled.png")
	imgprocess.GenerateImage(srcImg, tilerMC, 20, "tiled2.png")
	imgprocess.GenerateImage(srcImg, mosaicTiler, 50, "tiled3.png")
}
