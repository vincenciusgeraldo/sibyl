package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Review struct {
	Id         bson.ObjectId `json:"_id,omitempty" bson:"_id"`
	Repo       string        `json:"repo" bson:"repo"`
	PRNumber   int           `json:"pr_number" bson:"pr_number"`
	Requester  string        `json:"requester" bson:"requester"`
	Reviewers  []string      `json:"reviewers" bson:"reviewers"`
	ReviewedBy []string      `json:"reviewed_by" bson:"reviewed_by"`
	ApprovedBy []string      `json:"approved_by" bson:"approved_by"`
	Emergency  bool          `json:"emergency" bson:"emergency"`
	PRName     string		 `json:"pr_name" bson:"pr_name"`
	CreatedAt  time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at" bson:"updated_at"`
}
