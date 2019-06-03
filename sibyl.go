package sibyl

import (
	"github.com/vincenciusgeraldo/sibyl/pkg/interfaces"
	"github.com/vincenciusgeraldo/sibyl/pkg/models"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"github.com/google/go-github/github"
)

type Sibyl struct {
	bot interfaces.Bot
	rvw interfaces.Review
	usr interfaces.User
	git interfaces.PullRequestAPI
}

func NewSibyl(bot interfaces.Bot, rvw interfaces.Review, usr interfaces.User, pra interfaces.PullRequestAPI) *Sibyl {
	return &Sibyl{bot, rvw, usr, pra}
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
		chatId, _ := strconv.Atoi(usr.ChatId)
		rc := tb.User{ID: chatId}
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

func (s *Sibyl) GetUsers() ([]models.User, error) {
	return s.usr.GetUsers()
}

func (s *Sibyl) UpdateUser(user models.User) (models.User, error) {
	return s.usr.UpdateUser(user)
}

func (s *Sibyl) CreateReview(rvw models.Review) (models.Review, error) {
	return s.rvw.CreateReview(rvw)
}

func (s *Sibyl) GetReviewByRequester(usr string) ([]models.Review, error) {
	return s.rvw.GetReviewByRequester(usr)
}

func (s *Sibyl) GetReviewByReviewer(reviewer string, requester string) ([]models.Review, error) {
	return s.rvw.GetReviewByReviewer(reviewer, requester)
}

func (s *Sibyl) GetReviewBy(by interface{}) ([]models.Review, error) {
	return s.rvw.GetReviewBy(by)
}

func (s *Sibyl) GetApprovedReview() ([]models.Review, error) {
	return s.rvw.GetApprovedReview()
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

func (s *Sibyl) GetPullRequest(pr int, repo string) (*github.PullRequest, error) {
	return s.git.GetPullRequest(pr, repo)
}

func (s *Sibyl) GetPullRequestStatus(pr int, repo string) (models.PullRequestStatus, error) {
	return s.git.GetPullRequestStatus(pr, repo)
}