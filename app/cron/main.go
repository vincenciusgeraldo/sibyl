package main

import (
	"fmt"
	scheduler "github.com/robfig/cron"
	"os"
	"os/signal"
	"log"
	"github.com/vincenciusgeraldo/sibyl/pkg/database"
	"strconv"
	"github.com/vincenciusgeraldo/sibyl/pkg/repositories"
	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()
	go ClearReviews().Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

func ClearReviews() *scheduler.Cron{
	cron := scheduler.New()
	db, err := database.NewMongo(os.Getenv("MONGO_HOST"))
	if err != nil {
		log.Fatal(err)
	}

	rvw := repositories.NewReviewRepo(db)

	cron.AddFunc("@daily", func() {
		res, err := rvw.GetApprovedReview()
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("Failed to get approved reviews")
		}

		for _, r := range res {
			if err := rvw.DeleteReview(r.PRNumber, r.Repo); err != nil {
				fmt.Println("Failed to delete PR " + r.Repo + " " + strconv.Itoa(r.PRNumber))
			} else {
				fmt.Println("success to delete PR " + r.Repo + " " + strconv.Itoa(r.PRNumber))
			}
		}
	})
	return cron
}
