package main

import (
	"fmt"
	"github.com/subosito/gotenv"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	"time"
	"github.com/vincenciusgeraldo/sibyl"
	"github.com/vincenciusgeraldo/sibyl/pkg/database"
	"github.com/vincenciusgeraldo/sibyl/pkg/repositories"
	"github.com/vincenciusgeraldo/sibyl/pkg/handlers"
)

func main() {
	gotenv.Load()
	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 20 * time.Second},
		Reporter: func(err error) {
			fmt.Println(err)
		},
	})

	db, err := database.NewMongo(os.Getenv("MONGO_HOST"))

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Sibyl System is Up and Running")

	usr := repositories.NewUserRepo(db)
	rvw := repositories.NewReviewRepo(db)

	sbl := sibyl.NewSibyl(b, rvw, usr)

	ush := handlers.NewUserHandler(sbl)
	rvh := handlers.NewReviewHandler(sbl)
	adm := handlers.NewAdminHandler(sbl)

	b.Handle("/start", ush.Create)
	b.Handle("/help", ush.Help)
	b.Handle("/add", rvh.Create)
	b.Handle("/add_emergency", rvh.Create)
	b.Handle("/my_review", rvh.MyReview)
	b.Handle("/my_request", rvh.MyRequest)
	b.Handle("/reviewed", rvh.Reviewed)
	b.Handle("/approved", rvh.Approved)
	b.Handle("/done", rvh.Delete)
	b.Handle("/up", rvh.Up)
	b.Handle("/announce", adm.Announce)

	b.Start()
}
