package main

import (
	"github.com/nstehr/mosaicgen/db"
	"github.com/nstehr/mosaicgen/indexer/collect"
	"log"
	"os"
)

const (
	imgurClientID            = ""
	imgurClientSecret        = ""
	instagramClientID        = ""
	flickrKey                = ""
	flickrSecret             = ""
	twitterConsumerKey       = ""
	twitterConsumerSecret    = ""
	twitterAccessToken       = ""
	twitterAccessTokenSecret = ""
)

func main() {

	if len(os.Args) <= 1 {
		log.Fatal("not enough args, please pass in keyword")
	}
	keyword := os.Args[1]

	dbClient := db.NewMongoClient("localhost")
	defer dbClient.CloseConnection()

	instagramClient := collect.InstagramClient{ClientID: instagramClientID}
	imgurClient := collect.ImgurClient{ClientID: imgurClientID, ClientSecret: imgurClientSecret}
	flickrClient := collect.FlickrClient{Key: flickrKey, Secret: flickrSecret}
	twitterClient := collect.TwitterClient{ConsumerKey: twitterConsumerKey, ConsumerSecret: twitterConsumerSecret, AccessToken: twitterAccessToken, AccessTokenSecret: twitterAccessTokenSecret}
	sources := []collect.Source{instagramClient, imgurClient, flickrClient, twitterClient}

	collect.StartCollection(keyword, sources, dbClient)

}
