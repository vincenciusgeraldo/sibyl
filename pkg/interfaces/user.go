package interfaces

import "github.com/vincenciusgeraldo/sibyl/pkg/models"

type User interface {
	CreateUser(models.User) (models.User, error)
	GetUser(string) (models.User, error)
}
