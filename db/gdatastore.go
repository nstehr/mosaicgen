package db

import (
	"encoding/json"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/datastore"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
)

type GDatastoreClient struct {
	ctx context.Context
}

type gDataPhoto struct {
	Url      string `datastore:",noindex"`
	ThumbUrl string `datastore:",noindex"`
	Source   string
	Text     string `datastore:",noindex"`
	Tag      string
	AgvColor string `datastore:",noindex"`
}

func NewGDatastoreClientFromJSON(projectId string, jsonKeyPath string) *GDatastoreClient {
	jsonKey, err := ioutil.ReadFile(jsonKeyPath)
	if err != nil {
		log.Fatal(err)
	}
	conf, err := google.JWTConfigFromJSON(oauth2.NoContext, jsonKey, datastore.ScopeDatastore, datastore.ScopeUserEmail)

	if err != nil {
		log.Fatal(err)
	}

	ctx := cloud.NewContext(projectId, conf.Client(oauth2.NoContext, nil))

	return &GDatastoreClient{ctx: ctx}

}

func NewGDatastoreClientForComputeEngine(projectId string) *GDatastoreClient {
	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.ComputeTokenSource(""),
		},
	}
	ctx := cloud.NewContext(projectId, client)
	return &GDatastoreClient{ctx: ctx}
}

func (client *GDatastoreClient) SavePhoto(photo *Photo) {
	if client.ctx.Err() != nil {
		log.Fatal("Error saving photo: ")
	}
	gPhoto := gDataPhoto{Url: photo.Url, ThumbUrl: photo.ThumbUrl, Source: photo.Source, Text: photo.Text, Tag: photo.Tag}
	//datastore API can't take unsigned ints, so convert the whole struct to JSON
	//instead of doing int conversions
	avgColorJson, err := json.Marshal(photo.AvgColor)
	if err != nil {
		log.Println("error converting avg color")
	}
	gPhoto.AgvColor = string(avgColorJson)
	key := datastore.NewIncompleteKey(client.ctx, "Photo", nil)
	_, err = datastore.Put(client.ctx, key, &gPhoto)
	if err != nil {
		log.Printf("Error saving photo: %s error: %s\n", photo.Url, err)
	}
}

func (client *GDatastoreClient) CloseConnection() {
	//TODO: should be a way to use the client.ctx.Done() channel here
}

func (client *GDatastoreClient) GetPhotos(tag string, photos *[]Photo) {
	it := datastore.NewQuery("Photo").Filter("Tag = ", tag).Run(client.ctx)

	gPhoto := gDataPhoto{}
	hasData := true

	for hasData {
		_, err := it.Next(&gPhoto)
		if err != nil {
			hasData = false
		} else {
			avgColor := color.RGBA{}
			err := json.Unmarshal([]byte(gPhoto.AgvColor), &avgColor)
			if err != nil {
				log.Fatal("error unmarshalling avg color")
			}
			*photos = append(*photos, Photo{Url: gPhoto.Url, ThumbUrl: gPhoto.ThumbUrl, Source: gPhoto.Source, Text: gPhoto.Text, Tag: gPhoto.Tag, AvgColor: avgColor})

		}

	}

}
