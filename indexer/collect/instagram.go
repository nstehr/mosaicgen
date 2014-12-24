package collect

import (
	"github.com/gedex/go-instagram/instagram"
	"log"
	"time"
)

type InstagramClient struct {
	ClientID string
	api      *instagram.Client
}

func (client InstagramClient) Collect(searchTerm string) <-chan Photo {

	out := make(chan Photo)

	instagramClient := instagram.NewClient(nil)
	instagramClient.ClientID = client.ClientID
	client.api = instagramClient
	go client.getPictures(searchTerm, out)

	return out
}

func (client InstagramClient) getPictures(searchTerm string, ch chan Photo) {

	p := new(instagram.Parameters)
	moreData := true

	for moreData {
		media, next, err := client.api.Tags.RecentMedia(searchTerm, p)
		if err != nil {
			log.Fatal(err)
		}
		for _, m := range media {
			if m.Type == "image" {
				ph := Photo{}
				ph.Source = "instagram"
				ph.Tag = searchTerm
				if m.Caption != nil {
					ph.Text = m.Caption.Text
				}

				ph.ThumbUrl = m.Images.Thumbnail.URL
				ph.Url = m.Images.Thumbnail.URL
				ch <- ph
			}

		}
		if next.NextURL != "" {
			log.Println("waiting to make next instagram call")
			p.MaxID = next.NextMaxID
			time.Sleep(1 * time.Second)

		} else {
			moreData = false
		}
	}

	close(ch)
}
