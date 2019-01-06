package handlers

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/vincenciusgeraldo/sibyl"
	"github.com/vincenciusgeraldo/sibyl/pkg/models"
	"github.com/vincenciusgeraldo/sibyl/pkg/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
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
		msg = fmt.Sprintf("Antrian _review_ untuk PR `#%d %s` sudah ada. Silahkan coba `/up [repo] [pr_number]`", res[0].PRNumber, res[0].Repo)
		h.sbl.SendMessage(m.Chat, msg)
		return
	}

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
	}

	if _, err := h.sbl.CreateReview(review); err != nil {
		fmt.Println(err.Error())
		msg = fmt.Sprintf("Gagal meminta review untuk PR #%d %s.", prNumber, cmd[1])
		h.sbl.SendMessage(m.Chat, msg)
		return
	}

	if review.Emergency {
		tag = "\xE2\x80\xBC *Emergency* \xE2\x80\xBC\n"
	}

	if m.Sender.ID != int(m.Chat.ID) {
		msg = tag + fmt.Sprintf("Kak %s meminta review untuk PR `#%d %s` ke %s. Silahkan cek [disini](https://github.com/bukalapak/%s/pull/%d).", review.Requester, prNumber, cmd[1], cmd[3], cmd[1], prNumber)
		h.sbl.SendMessage(m.Chat, msg)
	} else {
		msg = tag + fmt.Sprintf("Antrian review untuk `PR #%d %s` berhasil dibuat dan diumumkan.", prNumber, cmd[1])
		h.sbl.SendMessage(m.Chat, msg)
	}

	msg = tag + fmt.Sprintf("%s meminta review untuk PR `#%d %s`. Silahkan cek [disini](https://github.com/bukalapak/%s/pull/%d).", review.Requester, prNumber, cmd[1], cmd[1], prNumber)
	h.sbl.BroadcastMessage(review.Reviewers, msg)

	return
}

func (h *Review) MyReview(m *tb.Message) {
	cmd := strings.Split(m.Text, " ")
	var res []models.Review
	var err error

	if len(cmd) == 1 {
		res, err = h.sbl.GetReviewByReviewer("@" + m.Sender.Username, "")
	} else if len(cmd) == 2 {
		res, err = h.sbl.GetReviewByReviewer("@" + m.Sender.Username, cmd[1])
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

	h.sbl.SendMessage(m.Chat, h.buildReviewList(res))
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

	h.sbl.SendMessage(m.Chat, h.buildReviewList(res))
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
		msg = fmt.Sprintf("Antrian _review_ untuk PR `#%d %s` tidak ditemukan.`", rev[0].PRNumber, rev[0].Repo)
		h.sbl.SendMessage(m.Chat, msg)
		return
	}

	res, err := h.sbl.Reviewed(pr, cmd[1], reviewer)
	if err != nil {
		h.sbl.SendMessage(m.Chat, fmt.Sprintf("Gagal memberikan _review_ untuk PR #%d %s", pr, cmd[1]))
		return
	}

	if m.Sender.ID != int(m.Chat.ID) {
		msg = fmt.Sprintf("Kak %s, PR `#%d %s` sudah di _review_ oleh %s. Silahkan cek [disini](https://github.com/bukalapak/%s/pull/%d).", res.Requester, pr, cmd[1], reviewer, cmd[1], pr)
		h.sbl.SendMessage(m.Chat, msg)
	} else {
		msg = fmt.Sprintf("_Review_ untuk PR `#%d %s` berhasil dan telah diumumkan.", pr, cmd[1])
		h.sbl.SendMessage(m.Chat, msg)
	}

	msg = fmt.Sprintf("PR `#%d %s` sudah di _review_ oleh %s. Silahkan cek [disini](https://github.com/bukalapak/%s/pull/%d).", pr, cmd[1], reviewer, cmd[1], pr)
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
		msg = fmt.Sprintf("Antrian _review_ untuk PR `#%d %s` tidak ditemukan.`", pr, cmd[1])
		h.sbl.SendMessage(m.Chat, msg)
		return
	}

	res, err := h.sbl.Approved(pr, cmd[1], reviewer)
	if err != nil {
		h.sbl.SendMessage(m.Chat, fmt.Sprintf("Gagal memberikan _approval_ untuk PR #%d %s", pr, cmd[1]))
		return
	}

	if m.Sender.ID != int(m.Chat.ID) {
		msg = fmt.Sprintf("Kak %s, PR `#%d %s` sudah di _approve_ oleh %s. Silahkan cek [disini](https://github.com/bukalapak/%s/pull/%d).", res.Requester, pr, cmd[1], reviewer, cmd[1], pr)
		h.sbl.SendMessage(m.Chat, msg)
	} else {
		msg = fmt.Sprintf("_Review_ untuk PR `#%d %s` berhasil dan telah diumumkan.", pr, cmd[1])
		h.sbl.SendMessage(m.Chat, msg)
	}

	msg = fmt.Sprintf("PR `#%d %s` sudah di _approve_ oleh %s. Silahkan cek [disini](https://github.com/bukalapak/%s/pull/%d).", pr, cmd[1], reviewer, cmd[1], pr)
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
		msg = fmt.Sprintf("Antrian _review_ untuk PR `#%d %s` tidak ditemukan.`", pr, cmd[1])
		h.sbl.SendMessage(m.Chat, msg)
		return
	}

	res, err := h.sbl.UpReview(pr, cmd[1])
	if err != nil {
		h.sbl.SendMessage(m.Chat, fmt.Sprintf("Gagal meminta _review_ untuk PR #%d %s", pr, cmd[1]))
		return
	}

	if m.Sender.ID != int(m.Chat.ID) {
		msg = fmt.Sprintf("Kak %s meminta review untuk PR `#%d %s` ke %s. Silahkan cek [disini](https://github.com/bukalapak/%s/pull/%d).", res.Requester, pr, cmd[1], strings.Join(res.Reviewers, " "), cmd[1], pr)
		h.sbl.SendMessage(m.Chat, msg)
	} else {
		msg = fmt.Sprintf("Antrian review untuk PR `#%d %s` berhasil diumumkan.", pr, cmd[1])
		h.sbl.SendMessage(m.Chat, msg)
	}

	msg = fmt.Sprintf("%s meminta review untuk PR `#%d %s`. Silahkan cek [disini](https://github.com/bukalapak/%s/pull/%d).", res.Requester, pr, cmd[1], cmd[1], pr)
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

func (h *Review) buildReviewList(data []models.Review) string {
	msg := ""
	for _, d := range data {
		if d.Emergency {
			msg += "\xE2\x80\xBC *Emergency* \xE2\x80\xBC\n" +
				fmt.Sprintf("Antrian PR `#%d %s` dari %s memerlukan _review_ secepatnya, Silahkan cek [disini](https://github.com/bukalapak/%s/pull/%d).\n", d.PRNumber, d.Repo, d.Requester, d.Repo, d.PRNumber)
		} else {
			msg += fmt.Sprintf("Antrian _review_ PR `#%d %s` dari %s. Silahkan cek [disini](https://github.com/bukalapak/%s/pull/%d).\n", d.PRNumber, d.Repo, d.Requester, d.Repo, d.PRNumber)
		}
		for _, r := range d.Reviewers {
			rst := ""
			if utils.ArrayInclude(d.ApprovedBy, r) {
				rst = "`approved` \xE2\x9C\x85"
			} else if utils.ArrayInclude(d.ReviewedBy, r) {
				rst = "`need changes` \xE2\x9D\x8C"
			} else {
				rst = "`not reviewed`"
			}
			msg += fmt.Sprintf("- %s - %s\n", r, rst)
		}
		msg += "\n"
	}
	return msg
}
