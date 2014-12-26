package collect

import (
	"bitbucket.org/liamstask/go-imgur/imgur"
	"github.com/nstehr/mosaicgen/db"
	"log"
	"time"
)

type ImgurClient struct {
	ClientID     string
	ClientSecret string
	api          *imgur.Client
}

const (
	apiSleepTime = 2 * time.Second
)

func (client ImgurClient) Collect(searchTerm string) <-chan db.Photo {

	out := make(chan db.Photo)

	imgurClient := imgur.NewClient(nil, client.ClientID, client.ClientSecret)
	client.api = imgurClient
	go client.getPictures(searchTerm, out)

	return out
}

func (client ImgurClient) getPictures(searchTerm string, ch chan db.Photo) {
	defer close(ch)
	moreData := true
	page := 1
	for moreData {
		results, err := client.api.Gallery.Search(searchTerm, "time", page)
		if err != nil {
			log.Printf("error retrieving pictures from imgur: %s\n", err)
			//increase the page number, and try again
			page++
			time.Sleep(apiSleepTime)
			continue
		}
		for _, r := range results {
			if !r.IsAlbum && !r.Animated {
				ph := db.Photo{}
				ph.Source = "imgur"
				ph.Tag = searchTerm
				ph.Text = r.Description
				ph.Url = r.Link

				ch <- ph
			}
		}

		if len(results) > 0 {
			page++
			log.Println("waiting to make next imgur call")
			time.Sleep(apiSleepTime)
		} else {
			moreData = false
		}
	}

}
