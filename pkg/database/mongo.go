package database

import (
	"github.com/globalsign/mgo"
	"os"
)

type Mongo struct {
	db *mgo.Database
}

func NewMongo(url string) (*mgo.Database, error) {
	session, err := mgo.Dial(url)

	if err != nil {
		return &mgo.Database{}, err
	}

	return session.DB(os.Getenv("MONGO_DATABASE")), nil
}
