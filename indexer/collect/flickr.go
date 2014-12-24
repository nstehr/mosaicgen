package collect

import (
	"github.com/manki/flickgo"
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

func (client FlickrClient) Collect(searchTerm string) <-chan Photo {

	out := make(chan Photo)
	flickrClient := flickgo.New(client.Key, client.Secret, http.DefaultClient)

	client.api = flickrClient

	go client.getPictures(searchTerm, out)

	return out
}

func (client FlickrClient) getPictures(searchTerm string, ch chan Photo) {
	flickrArgs := make(map[string]string)
	flickrArgs["tags"] = searchTerm
	page := 1
	moreData := true
	for moreData {
		flickrArgs["page"] = strconv.Itoa(page)
		resp, err := client.api.Search(flickrArgs)
		if err != nil {
			log.Fatal("error retrieving picture from flickr")
		}
		for _, photo := range resp.Photos {
			ph := Photo{}
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
			time.Sleep(1500 * time.Millisecond)
		} else {
			moreData = false
		}
	}
	close(ch)
}
