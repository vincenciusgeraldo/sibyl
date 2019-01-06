package repositories

import (
	"github.com/globalsign/mgo"
	"github.com/vincenciusgeraldo/sibyl/pkg/models"
	"github.com/globalsign/mgo/bson"
	"github.com/vincenciusgeraldo/sibyl/pkg/utils"
	"time"
)

type Review struct {
	db *mgo.Database
}

func NewReviewRepo(db *mgo.Database) *Review {
	return &Review{db}
}

func (r *Review) CreateReview(rvw models.Review) (models.Review, error) {
	rvw.CreatedAt = time.Now()
	rvw.UpdatedAt = time.Now()
	rvw.Id = bson.NewObjectId()

	if err := r.db.C("reviews").Insert(rvw); err != nil {
		return models.Review{}, err
	}

	return rvw, nil
}

func (r *Review) GetReviewByRequester(usr string) ([]models.Review, error){
	var res []models.Review
	q := map[string]interface{}{
		"requester": usr,
	}

	if err := r.db.C("reviews").Find(q).Sort("created_at").All(&res); err != nil {
		return []models.Review{}, err
	}

	return res, nil
}

func (r *Review) GetReviewByReviewer(reviewer string, requester string) ([]models.Review, error) {
	var res []models.Review
	q := bson.M{
		"reviewers": bson.M{"$in": []string{reviewer}},
		"approved_by": bson.M{"$nin": []string{reviewer}},
	}
	if requester != "" {
		q["requester"] = requester
	}

	if err := r.db.C("reviews").Find(q).Sort("created_at", "emergency").All(&res); err != nil {
		return []models.Review{}, err
	}

	return res, nil
}

func (r *Review) GetReviewBy(by interface{}) ([]models.Review, error) {
	var res []models.Review

	if err := r.db.C("reviews").Find(by).Sort("created_at", "emergency").All(&res); err != nil {
		return []models.Review{}, err
	}

	return res, nil
}

func (r *Review) Reviewed(pr int, repo string, usr string) (models.Review, error) {
	var rev []models.Review
	q := bson.M{
		"pr_number": pr,
		"repo": repo,
	}

	if err := r.db.C("reviews").Find(q).All(&rev); err != nil {
		return models.Review{},err
	}

	res := rev[0]
	res.ApprovedBy = utils.DeleteFromArray(usr, res.ApprovedBy)
	res.ReviewedBy = utils.UniqueArray(append(res.ReviewedBy, usr))
	res.UpdatedAt = time.Now()

	if err := r.db.C("reviews").Update(q, res); err != nil {
		return models.Review{},err
	}

	return res, nil
}

func (r *Review) Approved(pr int, repo string, usr string) (models.Review, error) {
	var rev []models.Review
	q := bson.M{
		"pr_number": pr,
		"repo": repo,
	}

	if err := r.db.C("reviews").Find(q).All(&rev); err != nil {
		return models.Review{},err
	}

	res := rev[0]
	res.ReviewedBy = utils.DeleteFromArray(usr, res.ReviewedBy)
	res.ApprovedBy = utils.UniqueArray(append(res.ApprovedBy, usr))
	res.UpdatedAt = time.Now()

	if err := r.db.C("reviews").Update(q, res); err != nil {
		return models.Review{},err
	}

	return res,nil
}

func (r *Review) UpReview(pr int, repo string) (models.Review, error) {
	var rev []models.Review
	q := bson.M{
		"pr_number": pr,
		"repo": repo,
	}

	if err := r.db.C("reviews").Find(q).All(&rev); err != nil {
		return models.Review{},err
	}

	res := rev[0]
	res.UpdatedAt = time.Now()

	if err := r.db.C("reviews").Update(q, res); err != nil {
		return models.Review{},err
	}

	return res, nil
}

func (r *Review) DeleteReview(pr int, repo string) error {
	q := bson.M{
		"pr_number": pr,
		"repo": repo,
	}

	if err := r.db.C("reviews").Remove(q); err != nil {
		return err
	}

	return nil
}