package db

type PhotoDB interface {
	SavePhoto(photo *Photo)
	GetPhotos(tag string, photos *[]Photo)
	CloseConnection()
}
