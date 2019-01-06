package handlers

import (
	"github.com/vincenciusgeraldo/sibyl/pkg/models"
	tb "gopkg.in/tucnak/telebot.v2"
	"github.com/vincenciusgeraldo/sibyl"
)

type User struct {
	sbl *sibyl.Sibyl
}

func NewUserHandler(sbl *sibyl.Sibyl) *User {
	return &User{sbl}
}

func (h *User) Create(m *tb.Message) {
	user := models.User{
		ChatId:   m.Sender.ID,
		Username: "@" + m.Sender.Username,
		Name:     m.Sender.FirstName + " " + m.Sender.LastName,
	}

	usr, err := h.sbl.GetUser(user.Username)
	if err != nil || usr.Username == "" {
		usr, err = h.sbl.CreateUser(user)
		if err != nil {
			h.sbl.SendMessage(m.Chat, "Gagal membuat user. Silahkan hubungi @vgeraldo untuk masalah terkait.")
			return
		}

		h.sbl.SendMessage(m.Chat, h.userHelper(usr.Name))
		return
	} else {
		h.sbl.SendMessage(m.Chat, "Kamu sudah terdaftar sebagai pengguna Sibyl. Selamat menggunakan Sibyl :)")
		return
	}
}

func (h *User) Help(m *tb.Message) {
	h.sbl.SendMessage(m.Chat, h.userHelper(""))
	return
}

func (h *User) userHelper(usr string) string {
	grt := ""
	if usr != "" {
		grt = "Hi, " + usr + "\n" +
			"Terimakasih sudah menggunakan Sibyl.\n\n"
	}

	return grt + "*Apa itu Sibyl?*\n" +
		"Sybil merupakan _tools_ untuk mempermudah antrian _review_.\n\n" +

		"*Apa yang bisa dilakukan Sibyl?*\n" +
		"1. Menambahkan antrian _review_ baru.\n" +
		"`/add [repo] [pr_number] [reviewers]`\n" +
		"2. Melihat antrian _review_ mu.\n" +
		"`/my_request`\n" +
		"3. Melihat antrian _review_ yang butuh review kamu.\n" +
		"`/my_review [requester (optional)]`\n" +
		"4. Menandai PR sudah di _review_\n" +
		"`/reviewed [repo] [pr_number]`\n" +
		"5. Menandai PR sudah di _approve_\n" +
		"`/approved [repo] [pr_number]`\n" +
		"6. Menghapus antrian _review_\n" +
		"`/done [repo] [pr_number]`\n" +
		"7. Meminta _review_ kembali\n" +
		"`/up [repo] [pr_number]`\n" +
		"8. Menambahkan antrian _emergency review_\n" +
		"`/add_emergency [repo] [pr_number] [reviewers]`\n"
}

func (h *User) userExist(usr string) string {
	return "Hi, " + usr + "\n" +
		"Dirimu sudah terdaftar sebagai pengguna Sibyl. Selamat menggunakan :)"
}
