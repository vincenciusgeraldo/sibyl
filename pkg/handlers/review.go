package handlers

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/vincenciusgeraldo/sibyl"
	"github.com/vincenciusgeraldo/sibyl/pkg/models"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
	"github.com/vincenciusgeraldo/sibyl/pkg/common"
)

type Review struct {
	sbl *sibyl.Sibyl
}

func NewReviewHandler(sbl *sibyl.Sibyl) *Review {
	return &Review{sbl}
}

func (h *Review) Create(m *tb.Message) {
	msg, tag := "", ""
	cmd := strings.SplitN(m.Text, " ", 4)
	if len(cmd) < 4 {
		h.sbl.SendMessage(m.Chat, "Silahkan coba `/add [repo] [pr_number] [list_reviewer]`")
		return
	}

	if usr, err := h.sbl.GetUser("@" + m.Sender.Username); err != nil {
		h.sbl.SendMessage(m.Chat, "Terjadi kesalahan pada Sibyl. Silahkan hubungi @vgeraldo untuk masalah terkait.")
		return
	} else if usr.Username == "" {
		h.sbl.SendMessage(m.Chat, "Kamu belum terdaftar sebagai pengguna Sibyl. Silahkan chat /start ke Sibyl Bot untuk mendapatkan notifikasi.")
	}

	prNumber, err := strconv.Atoi(cmd[2])
	if err != nil {
		h.sbl.SendMessage(m.Chat, "Silahkan coba `/add [repo] [pr_number] [list_reviewer]`")
		return
	}

	res, err := h.sbl.GetReviewBy(bson.M{"repo": cmd[1], "pr_number": prNumber})
	if err != nil {
		fmt.Println(err.Error())
		h.sbl.SendMessage(m.Chat, "Terjadi kesalahan pada Sibyl. Silahkan hubungi @vgeraldo untuk masalah terkait.")
		return
	}
	if len(res) == 1 {
		msg = fmt.Sprintf("Antrian _review_ untuk PR `%s#%d` sudah ada. Silahkan coba `/up [repo] [pr_number]`", res[0].Repo, res[0].PRNumber)
		h.sbl.SendMessage(m.Chat, msg)
		return
	}

	pr, _ := h.sbl.GetPullRequest(prNumber, cmd[1])

	emg := false
	if cmd[0] == "/add_emergency" {
		emg = true
	}
	review := models.Review{
		Repo:      cmd[1],
		PRNumber:  prNumber,
		Requester: "@" + m.Sender.Username,
		Reviewers: strings.Split(cmd[3], " "),
		Emergency: emg,
		PRName: pr.GetTitle(),
	}

	if _, err := h.sbl.CreateReview(review); err != nil {
		fmt.Println(err.Error())
		msg = fmt.Sprintf("Gagal meminta review untuk PR %s#%d.", cmd[1], prNumber)
		h.sbl.SendMessage(m.Chat, msg)
		return
	}

	if review.Emergency {
		tag = "\xE2\x80\xBC *Emergency* \xE2\x80\xBC\n"
	}

	if m.Sender.ID != int(m.Chat.ID) {
		msg = tag + fmt.Sprintf("Kak %s meminta review untuk PR `%s#%d` - `%s` ke %s. PR-nya bisa dilihat [disini](https://github.com/bukalapak/%s/pull/%d).", review.Requester, cmd[1], prNumber, review.PRName, cmd[3], cmd[1], prNumber)
		h.sbl.SendMessage(m.Chat, msg)
	} else {
		msg = tag + fmt.Sprintf("Antrian review untuk `PR %s#%d` berhasil dibuat dan diumumkan.", cmd[1], prNumber)
		h.sbl.SendMessage(m.Chat, msg)
	}

	msg = tag + fmt.Sprintf("%s meminta review untuk PR `%s#%d` - `%s`. PR-nya bisa dilihat [disini](https://github.com/bukalapak/%s/pull/%d).", review.Requester, cmd[1], prNumber, review.PRName, cmd[1], prNumber)
	h.sbl.BroadcastMessage(review.Reviewers, msg)

	return
}

func (h *Review) MyReview(m *tb.Message) {
	cmd := strings.Split(m.Text, " ")
	var res []models.Review
	var err error
	if len(cmd) == 1 {
		res, err = h.sbl.GetReviewByReviewer("@"+m.Sender.Username, "")
	} else if len(cmd) == 2 {
		res, err = h.sbl.GetReviewByReviewer("@"+m.Sender.Username, cmd[1])
	} else {
		h.sbl.SendMessage(m.Chat, "Silahkan coba `/my_review [ requester (optional)]`")
		return
	}

	if err != nil {
		fmt.Println(err.Error())
		h.sbl.SendMessage(m.Chat, "Terjadi kesalahan pada Sibyl. Silahkan hubungi @vgeraldo untuk masalah terkait.")
		return
	}

	if len(res) == 0 {
		h.sbl.SendMessage(m.Chat, "Kamu belum memiliki antrian _review_.")
		return
	}

	h.buildReviewList(res, m)
	return
}

func (h *Review) MyRequest(m *tb.Message) {
	res, err := h.sbl.GetReviewByRequester("@" + m.Sender.Username)
	if err != nil {
		fmt.Println(err.Error())
		h.sbl.SendMessage(m.Chat, "Terjadi kesalahan pada Sibyl. Silahkan hubungi @vgeraldo untuk masalah terkait.")
		return
	}

	if len(res) == 0 {
		h.sbl.SendMessage(m.Chat, "Kamu belum memiliki permintaan _review_.")
		return
	}

	h.buildReviewList(res, m)
	return
}

func (h *Review) Reviewed(m *tb.Message) {
	msg := ""
	reviewer := "@" + m.Sender.Username
	cmd := strings.SplitN(m.Text, " ", 3)
	if len(cmd) < 3 {
		h.sbl.SendMessage(m.Chat, "Silahkan coba `/reviewed [repo] [pr_number]`")
		return
	}

	pr, err := strconv.Atoi(cmd[2])
	if err != nil {
		h.sbl.SendMessage(m.Chat, "Silahkan coba `/reviewed [repo] [pr_number]`")
		return
	}

	rev, err := h.sbl.GetReviewBy(bson.M{"repo": cmd[1], "pr_number": pr})
	if err != nil {
		fmt.Println(err.Error())
		h.sbl.SendMessage(m.Chat, "Terjadi kesalahan pada Sibyl. Silahkan hubungi @vgeraldo untuk masalah terkait.")
		return
	}
	if len(rev) == 0 {
		msg = fmt.Sprintf("Antrian _review_ untuk PR `%s#%d` tidak ditemukan.`", rev[0].Repo, rev[0].PRNumber)
		h.sbl.SendMessage(m.Chat, msg)
		return
	}

	res, err := h.sbl.Reviewed(pr, cmd[1], reviewer)
	if err != nil {
		h.sbl.SendMessage(m.Chat, fmt.Sprintf("Gagal memberikan _review_ untuk PR %s#%d", cmd[1], pr))
		return
	}

	if m.Sender.ID != int(m.Chat.ID) {
		msg = fmt.Sprintf("Kak %s, PR `%s#%d` - `%s` sudah di _review_ oleh %s. PR-nya bisa dilihat [disini](https://github.com/bukalapak/%s/pull/%d).", res.Requester, cmd[1], pr, res.PRName, reviewer, cmd[1], pr)
		h.sbl.SendMessage(m.Chat, msg)
	} else {
		msg = fmt.Sprintf("_Review_ untuk PR `%s#%d` berhasil dan telah diumumkan.", cmd[1], pr)
		h.sbl.SendMessage(m.Chat, msg)
	}

	msg = fmt.Sprintf("PR `%s#%d` - `%s` sudah di _review_ oleh %s. PR-nya bisa dilihat [disini](https://github.com/bukalapak/%s/pull/%d).", cmd[1], pr, reviewer, res.PRName, cmd[1], pr)
	h.sbl.BroadcastMessage([]string{res.Requester}, msg)
	return
}

func (h *Review) Approved(m *tb.Message) {
	msg := ""
	reviewer := "@" + m.Sender.Username
	cmd := strings.SplitN(m.Text, " ", 3)

	if len(cmd) < 3 {
		h.sbl.SendMessage(m.Chat, "Silahkan coba `/approved [repo] [pr_number]`")
		return
	}

	pr, err := strconv.Atoi(cmd[2])
	if err != nil {
		h.sbl.SendMessage(m.Chat, "Silahkan coba `/approved [repo] [pr_number]`")
		return
	}

	rev, err := h.sbl.GetReviewBy(bson.M{"repo": cmd[1], "pr_number": pr})
	if err != nil {
		fmt.Println(err.Error())
		h.sbl.SendMessage(m.Chat, "Terjadi kesalahan pada Sibyl. Silahkan hubungi @vgeraldo untuk masalah terkait.")
		return
	}

	if len(rev) == 0 {
		msg = fmt.Sprintf("Antrian _review_ untuk PR `%s#%d` tidak ditemukan.`", cmd[1], pr)
		h.sbl.SendMessage(m.Chat, msg)
		return
	}

	res, err := h.sbl.Approved(pr, cmd[1], reviewer)
	if err != nil {
		h.sbl.SendMessage(m.Chat, fmt.Sprintf("Gagal memberikan _approval_ untuk PR %s#%d", cmd[1], pr))
		return
	}

	usr, err := h.sbl.GetUser(reviewer)
	if err != nil {
		usr = models.User{}
	}

	if usr.Role == 1 {
		h.sbl.DeleteReview(pr, cmd[1])
	}

	if m.Sender.ID != int(m.Chat.ID) {
		msg = fmt.Sprintf("Kak %s, PR `%s#%d` - `%s` sudah di _approve_ oleh %s. PR-nya bisa dilihat [disini](https://github.com/bukalapak/%s/pull/%d).", res.Requester, cmd[1], pr, res.PRName, reviewer, cmd[1], pr)
		h.sbl.SendMessage(m.Chat, msg)
	} else {
		msg = fmt.Sprintf("_Review_ untuk PR `%s#%d` berhasil dan telah diumumkan.", cmd[1], pr)
		h.sbl.SendMessage(m.Chat, msg)
	}

	msg = fmt.Sprintf("PR `%s#%d` - `%s` sudah di _approve_ oleh %s. PR-nya bisa dilihat [disini](https://github.com/bukalapak/%s/pull/%d).", cmd[1], pr, res.PRName, reviewer, cmd[1], pr)
	h.sbl.BroadcastMessage([]string{res.Requester}, msg)
	return
}

func (h *Review) Up(m *tb.Message) {
	msg := ""
	cmd := strings.SplitN(m.Text, " ", 3)
	if len(cmd) < 3 {
		h.sbl.SendMessage(m.Chat, "Silahkan coba `/up [repo] [pr_number]`")
		return
	}

	if usr, err := h.sbl.GetUser("@" + m.Sender.Username); err != nil {
		h.sbl.SendMessage(m.Chat, "Terjadi kesalahan pada Sibyl. Silahkan hubungi @vgeraldo untuk masalah terkait.")
		return
	} else if usr.Username == "" {
		h.sbl.SendMessage(m.Chat, "Kamu belum terdaftar sebagai pengguna Sibyl. Silahkan chat /start ke Sibyl Bot untuk mendapatkan notifikasi.")
	}

	pr, err := strconv.Atoi(cmd[2])
	if err != nil {
		h.sbl.SendMessage(m.Chat, "Silahkan coba `/up [repo] [pr_number]`")
		return
	}

	rev, err := h.sbl.GetReviewBy(bson.M{"repo": cmd[1], "pr_number": pr})
	if err != nil {
		fmt.Println(err.Error())
		h.sbl.SendMessage(m.Chat, "Terjadi kesalahan pada Sibyl. Silahkan hubungi @vgeraldo untuk masalah terkait.")
		return
	}
	if len(rev) == 0 {
		msg = fmt.Sprintf("Antrian _review_ untuk PR `%s#%d` tidak ditemukan.`", cmd[1], pr)
		h.sbl.SendMessage(m.Chat, msg)
		return
	}

	res, err := h.sbl.UpReview(pr, cmd[1])
	if err != nil {
		h.sbl.SendMessage(m.Chat, fmt.Sprintf("Gagal meminta _review_ untuk PR %s#%d", cmd[1], pr))
		return
	}

	if m.Sender.ID != int(m.Chat.ID) {
		msg = fmt.Sprintf("Kak %s meminta review untuk PR `%s#%d` - `%s` ke %s. PR-nya bisa dilihat [disini](https://github.com/bukalapak/%s/pull/%d).", res.Requester, cmd[1], pr, res.PRName, strings.Join(res.Reviewers, " "), cmd[1], pr)
		h.sbl.SendMessage(m.Chat, msg)
	} else {
		msg = fmt.Sprintf("Antrian review untuk PR `%s#%d` berhasil diumumkan.", cmd[1], pr)
		h.sbl.SendMessage(m.Chat, msg)
	}

	msg = fmt.Sprintf("%s meminta review untuk PR `%s#%d` - `%s`. PR-nya bisa dilihat [disini](https://github.com/bukalapak/%s/pull/%d).", res.Requester, cmd[1], pr, res.PRName, cmd[1], pr)
	h.sbl.BroadcastMessage(res.Reviewers, msg)

	return
}

func (h *Review) Delete(m *tb.Message) {
	msg := ""
	cmd := strings.SplitN(m.Text, " ", 3)
	if len(cmd) < 3 {
		h.sbl.SendMessage(m.Chat, "Silahkan coba `/done [repo] [pr_number]`")
		return
	}

	pr, err := strconv.Atoi(cmd[2])
	if err != nil {
		h.sbl.SendMessage(m.Chat, "Silahkan coba `/done [repo] [pr_number]`")
		return
	}

	rev, err := h.sbl.GetReviewBy(bson.M{"repo": cmd[1], "pr_number": pr})
	if err != nil {
		fmt.Println(err.Error())
		h.sbl.SendMessage(m.Chat, "Terjadi kesalahan pada Sibyl. Silahkan hubungi @vgeraldo untuk masalah terkait.")
		return
	}
	if len(rev) == 0 {
		msg = fmt.Sprintf("Antrian _review_ untuk PR `#%d %s` tidak ditemukan.`", pr, cmd[1])
		h.sbl.SendMessage(m.Chat, msg)
		return
	}

	err = h.sbl.DeleteReview(pr, cmd[1])
	if err != nil {
		h.sbl.SendMessage(m.Chat, fmt.Sprintf("Gagal menghapus antrian _review_ untuk PR #%d %s", pr, cmd[1]))
		return
	}

	msg = fmt.Sprintf("Antrian review untuk PR `#%d %s` berhasil dihapus.", pr, cmd[1])
	h.sbl.SendMessage(m.Chat, msg)

	return
}

func (h *Review) buildReviewList(data []models.Review, m *tb.Message) {
	msg := "Mohon di tunggu ya.. Sibyl sedang mengambil daftar permintaan review dan status PR nya.."
	h.sbl.SendMessage(m.Chat, msg)
	for _, d := range data {
		check, _ := h.sbl.GetPullRequestStatus(d.PRNumber, d.Repo)
		h.sbl.SendMessage(m.Chat, common.BuildReview(d, check))
	}
	msg = "Urutan prioritas review tertinggi dimulai dari paling bawah. Mohon _reviewer_ untuk memberikan _review_ dimulai dari antrian bawah."
	h.sbl.SendMessage(m.Chat, msg)
}
