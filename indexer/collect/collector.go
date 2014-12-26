package collect

import (
	"github.com/nstehr/mosaicgen/db"
	"github.com/nstehr/mosaicgen/imgprocess"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"sync"
)

type Source interface {
	Collect(searchTerm string) <-chan db.Photo
}

func StartCollection(searchTerm string, sources []Source, photoDB db.PhotoDB) {
	var wg sync.WaitGroup
	for _, source := range sources {
		//for each source, also create a channel for saving to the DB
		out := make(chan db.Photo)
		wg.Add(1)
		//for each source, launch a goroutine to process the images the source returns on it's channel
		go func(s Source, ch chan db.Photo) {
			defer wg.Done()
			for p := range s.Collect(searchTerm) {
				processPhoto(&p)
				ch <- p
			}
			close(ch)
		}(source, out)
		//for each source, launch a goroutine to read a channel and save the processed image to the DB
		wg.Add(1)
		go func(ch chan db.Photo) {
			defer wg.Done()
			for p := range ch {
				persistPhoto(&p, photoDB)
			}
		}(out)
	}

	//wait until all sources have closed their channels (no more pics)
	wg.Wait()
}

func processPhoto(photo *db.Photo) {
	var url string
	//use the thumb/small image for processing
	//if available
	if photo.ThumbUrl != "" {
		url = photo.ThumbUrl
	} else {
		url = photo.Url
	}
	log.Println("processing: " + url)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error downloading image from: " + url)
		return
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Println("Error decoding image from: " + url)
		return
	}
	r, g, b := imgprocess.GetAvgColor(img)
	photo.AvgColor = color.RGBA{r, g, b, 255}

}

func persistPhoto(photo *db.Photo, photoDB db.PhotoDB) {
	photoDB.SavePhoto(photo)
}
