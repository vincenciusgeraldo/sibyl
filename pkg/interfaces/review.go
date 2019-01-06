package interfaces

import (
	"github.com/vincenciusgeraldo/sibyl/pkg/models"
)

type Review interface {
	CreateReview(models.Review) (models.Review, error)
	GetReviewByRequester(string) ([]models.Review, error)
	GetReviewByReviewer(string, string) ([]models.Review, error)
	GetReviewBy(interface{}) ([]models.Review, error)
	Reviewed(int, string, string) (models.Review, error)
	Approved(int, string, string) (models.Review, error)
	UpReview(int, string) (models.Review, error)
	DeleteReview(int, string) error
}
