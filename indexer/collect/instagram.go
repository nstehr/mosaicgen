package collect

import (
	"github.com/gedex/go-instagram/instagram"
	"github.com/nstehr/mosaicgen/db"
	"log"
	"time"
)

type InstagramClient struct {
	ClientID string
	api      *instagram.Client
}

const (
	instagramAPISleepTime = 1 * time.Second
)

func (client InstagramClient) Collect(searchTerm string) <-chan db.Photo {

	out := make(chan db.Photo)

	instagramClient := instagram.NewClient(nil)
	instagramClient.ClientID = client.ClientID
	client.api = instagramClient
	go client.getPictures(searchTerm, out)

	return out
}

func (client InstagramClient) getPictures(searchTerm string, ch chan db.Photo) {
	defer close(ch)
	p := new(instagram.Parameters)
	moreData := true

	for moreData {
		media, next, err := client.api.Tags.RecentMedia(searchTerm, p)
		if err != nil {
			log.Printf("error retrieving pictures from instagram: %s\n", err)
			return
		}
		for _, m := range media {
			if m.Type == "image" {
				ph := db.Photo{}
				ph.Source = "instagram"
				ph.Tag = searchTerm
				if m.Caption != nil {
					ph.Text = m.Caption.Text
				}

				ph.ThumbUrl = m.Images.Thumbnail.URL
				ph.Url = m.Images.StandardResolution.URL
				ch <- ph
			}

		}
		if next.NextURL != "" {
			log.Println("waiting to make next instagram call")
			p.MaxID = next.NextMaxID
			time.Sleep(instagramAPISleepTime)

		} else {
			moreData = false
		}
	}
}
