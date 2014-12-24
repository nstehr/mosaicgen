package collect

import (
	"bitbucket.org/liamstask/go-imgur/imgur"
	"log"
	"time"
)

type ImgurClient struct {
	ClientID     string
	ClientSecret string
	api          *imgur.Client
}

func (client ImgurClient) Collect(searchTerm string) <-chan Photo {

	out := make(chan Photo)

	imgurClient := imgur.NewClient(nil, client.ClientID, client.ClientSecret)
	client.api = imgurClient
	go client.getPictures(searchTerm, out)

	return out
}

func (client ImgurClient) getPictures(searchTerm string, ch chan Photo) {

	moreData := true
	page := 1
	for moreData {
		results, err := client.api.Gallery.Search(searchTerm, "time", page)
		if err != nil {
			log.Fatal(err)
		}
		for _, r := range results {
			if !r.IsAlbum && !r.Animated {
				ph := Photo{}
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
			time.Sleep(2 * time.Second)
		} else {
			moreData = false
		}
	}

	close(ch)
}
