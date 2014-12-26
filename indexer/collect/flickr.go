package collect

import (
	"github.com/manki/flickgo"
	"github.com/nstehr/mosaicgen/db"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FlickrClient struct {
	Key    string
	Secret string
	api    *flickgo.Client
}

const (
	apiSleepTime = 1500 * time.Millisecond
)

func (client FlickrClient) Collect(searchTerm string) <-chan db.Photo {

	out := make(chan db.Photo)
	flickrClient := flickgo.New(client.Key, client.Secret, http.DefaultClient)

	client.api = flickrClient

	go client.getPictures(searchTerm, out)

	return out
}

func (client FlickrClient) getPictures(searchTerm string, ch chan db.Photo) {
	defer close(ch)
	flickrArgs := make(map[string]string)
	flickrArgs["tags"] = searchTerm
	page := 1
	moreData := true
	for moreData {
		flickrArgs["page"] = strconv.Itoa(page)
		resp, err := client.api.Search(flickrArgs)
		if err != nil {
			log.Printf("error retrieving pictures from flickr: %s\n", err)
			page++
			time.Sleep(apiSleepTime)
			continue
		}
		for _, photo := range resp.Photos {
			ph := db.Photo{}
			//what size should I use?
			ph.ThumbUrl = photo.URL(flickgo.SizeSmallSquare)
			ph.Url = photo.URL(flickgo.SizeMedium500)
			ph.Source = "flickr"
			ph.Tag = searchTerm
			ph.Text = photo.Title

			ch <- ph
		}
		if resp.Page < resp.Pages {
			page++
			log.Println("waiting to make next flickr call")
			time.Sleep(apiSleepTime)
		} else {
			moreData = false
		}
	}
}
