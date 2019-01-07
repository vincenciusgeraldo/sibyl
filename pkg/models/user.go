package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type User struct {
	Id        bson.ObjectId `json:"_id,omitempty" bson:"_id"`
	Username  string        `json:"username" bson:"username"`
	ChatId    string        `json:"chat_id" bson:"chat_id"`
	Name      string        `json:"name" bson:"name"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}
