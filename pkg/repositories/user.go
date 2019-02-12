package repositories

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/vincenciusgeraldo/sibyl/pkg/models"
	"time"
)

type User struct {
	db *mgo.Database
}

func NewUserRepo(db *mgo.Database) *User {
	return &User{db}
}

func (r *User) CreateUser(usr models.User) (models.User, error) {
	usr.Id = bson.NewObjectId()
	usr.CreatedAt = time.Now()
	usr.UpdatedAt = time.Now()

	if err := r.db.C("users").Insert(&usr); err != nil {
		return models.User{}, err
	}

	return usr, nil
}

func (r *User) UpdateUser(usr models.User) (models.User, error) {
	usr.UpdatedAt = time.Now()

	if err := r.db.C("users").UpdateId(usr.Id, usr); err != nil {
		return models.User{}, err
	}

	return usr, nil
}

func (r *User) GetUsers() ([]models.User, error) {
	var res []models.User

	if err := r.db.C("users").Find(bson.M{}).All(&res); err != nil {
		return []models.User{}, err
	}

	return res, nil
}

func (r *User) GetUser(usr string) (models.User, error) {
	var res []models.User
	q := map[string]interface{}{
		"username": usr,
	}

	if err := r.db.C("users").Find(q).All(&res); err != nil {
		return models.User{}, err
	}

	if len(res) == 0 {
		return models.User{}, nil
	}

	return res[0], nil
}
