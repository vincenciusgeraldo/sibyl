package sibyl

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"github.com/vincenciusgeraldo/sibyl/pkg/interfaces"
	"github.com/vincenciusgeraldo/sibyl/pkg/models"
)

type Sibyl struct {
	bot interfaces.Bot
	rvw interfaces.Review
	usr interfaces.User
}

func NewSibyl(bot interfaces.Bot, rvw interfaces.Review, usr interfaces.User) *Sibyl {
	return &Sibyl{bot, rvw, usr}
}

func (s *Sibyl) SendMessage(to tb.Recipient, m string) (*tb.Message, error) {
	return s.bot.Send(to, m, tb.ModeMarkdown)
}

func (s *Sibyl) BroadcastMessage(to []string, m string) error {
	for _, t := range to {
		usr, err := s.GetUser(t)
		if err != nil {
			return err
		}
		rc := tb.User{ID: usr.ChatId,}
		s.SendMessage(&rc, m)
	}

	return nil
}

func (s *Sibyl) CreateUser(user models.User) (models.User, error) {
	return s.usr.CreateUser(user)
}

func (s *Sibyl) GetUser(user string) (models.User, error) {
	return s.usr.GetUser(user)
}

func (s *Sibyl) CreateReview(rvw models.Review) (models.Review, error) {
	return s.rvw.CreateReview(rvw)
}

func (s *Sibyl) GetReviewByRequester(usr string) ([]models.Review, error){
	return s.rvw.GetReviewByRequester(usr)
}

func (s *Sibyl) GetReviewByReviewer(usr string) ([]models.Review, error) {
	return s.rvw.GetReviewByReviewer(usr)
}

func (s *Sibyl) GetReviewBy(by interface{}) ([]models.Review, error) {
	return s.rvw.GetReviewBy(by)
}

func (s *Sibyl) Reviewed(pr int, repo string, usr string) (models.Review, error) {
	return s.rvw.Reviewed(pr, repo, usr)
}

func (s *Sibyl) Approved(pr int, repo string, usr string) (models.Review, error) {
	return s.rvw.Approved(pr, repo, usr)
}

func (s *Sibyl) UpReview(pr int, repo string) (models.Review, error) {
	return s.rvw.UpReview(pr, repo)
}

func (s *Sibyl) DeleteReview(pr int, repo string) error {
	return s.rvw.DeleteReview(pr, repo)
}