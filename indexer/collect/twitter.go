package collect

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/nstehr/mosaicgen/db"
	"log"
	"net/url"
	"time"
)

type TwitterClient struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
	api               *anaconda.TwitterApi
}

const (
	twitterAPISleepTime = 2500 * time.Millisecond
)

func (client TwitterClient) Collect(searchTerm string) <-chan db.Photo {

	out := make(chan db.Photo)

	anaconda.SetConsumerKey(client.ConsumerKey)
	anaconda.SetConsumerSecret(client.ConsumerSecret)
	api := anaconda.NewTwitterApi(client.AccessToken, client.AccessTokenSecret)
	client.api = api
	go client.getPictures(searchTerm, out)
	return out
}

func (client TwitterClient) getPictures(searchTerm string, ch chan db.Photo) {
	defer close(ch)
	moreData := true
	v := url.Values{}
	maxId := int64(0)
	maxIdStr := ""
	for moreData {
		results, err := client.api.GetSearch(searchTerm, v)
		if err != nil {
			log.Printf("error retrieving pictures from twitter: %s\n", err)
		}
		for _, tweet := range results {
			if len(tweet.Entities.Media) > 0 {
				for _, media := range tweet.Entities.Media {
					if media.Type == "photo" {
						ph := db.Photo{}
						ph.Source = "twitter"
						ph.Url = media.Media_url
						ph.Tag = searchTerm
						ph.Text = tweet.Text

						ch <- ph
					}

				}

			}
			if tweet.Id > maxId {
				maxId = tweet.Id
				maxIdStr = tweet.IdStr
			}
		}
		if len(results) > 0 {
			log.Println("waiting to make next twitter call")
			v.Set("since_id", maxIdStr)
			//the anaconda api lib says it supports a delay via
			//api.SetDelay, but it panics when I use it.  It also
			//says it will handle rate limit errors, but I feel a bit
			//safer doing it myself....
			time.Sleep(twitterAPISleepTime)
		} else {
			moreData = false
		}
	}
}
