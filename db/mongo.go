package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type MongoDBClient struct {
	session *mgo.Session
}

const (
	database   = "photodb"
	collection = "photos"
)

func NewMongoClient(server string) *MongoDBClient {
	session, err := mgo.Dial(server)
	if err != nil {
		log.Fatal(err)
	}
	return &MongoDBClient{session: session}
}

func (client *MongoDBClient) SavePhoto(photo *Photo) {
	if client.session == nil {
		log.Fatal("connection to database closed!")
	}
	//could be called from a goroutine, so copy the session
	//to get its own temporary connection
	sessionCopy := client.session.Copy()
	defer sessionCopy.Close()

	sessionCopy.DB(database).C(collection).Insert(photo)
}

func (client *MongoDBClient) CloseConnection() {
	client.session.Close()
	client.session = nil
}

func (client *MongoDBClient) GetPhotos(tag string, photos *[]Photo) {
	if client.session == nil {
		log.Fatal("connection to database closed!")
	}
	err := client.session.DB(database).C(collection).Find(bson.M{"tag": tag}).All(photos)
	if err != nil {
		log.Printf("error getting photos : %s\n", err)
		return
	}
}
