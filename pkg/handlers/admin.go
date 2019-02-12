package handlers

import (
	"github.com/vincenciusgeraldo/sibyl"
	"github.com/vincenciusgeraldo/sibyl/pkg/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

type Admin struct {
	sbl *sibyl.Sibyl
}

func NewAdminHandler(sbl *sibyl.Sibyl) *Admin {
	return &Admin{sbl}
}

func (h *Admin) Announce(m *tb.Message) {
	msg := m.Payload
	if m.Sender.Username != "vgeraldo" {
		h.sbl.SendMessage(m.Chat, "Kamu tidak mempunya izin untuk menggunakan fitur ini.")
		return
	}

	usrs, err := h.sbl.GetUsers()
	if err != nil {
		h.sbl.SendMessage(m.Chat, "Terjadi kesalahan pada Sibyl. Silahkan hubungi @vgeraldo untuk masalah terkait.")
		return
	}

	arr := []string{}
	for _, usr := range usrs {
		if usr.Username != "vgeraldo" {
			arr = append(arr, usr.Username)
		}
	}
	arr = utils.UniqueArray(arr)

	h.sbl.BroadcastMessage(arr, msg)
	return
}

func (h *Admin) SetRole(m *tb.Message) {
	msg := strings.Split(m.Payload, " ")
	if m.Sender.Username != "vgeraldo" {
		h.sbl.SendMessage(m.Chat, "Kamu tidak mempunya izin untuk menggunakan fitur ini. Silahkan hubungi @vgeraldo untuk maminta akses.")
		return
	}

	usr, err := h.sbl.GetUser(msg[0])
	if err != nil {
		h.sbl.SendMessage(m.Chat, "Terjadi kesalahan pada Sibyl. Silahkan hubungi @vgeraldo untuk masalah terkait.")
		return
	}

	usr.SetRole(msg[1])

	if _, err = h.sbl.UpdateUser(usr); err != nil {
		h.sbl.SendMessage(m.Chat, "Gagal mengganti role " + msg[0] + " menjadi `" + msg[1] + "`")
		return
	}

	h.sbl.SendMessage(m.Chat, "Berhasil mengganti role " + msg[0] + " menjadi `" + msg[1] + "`")
	return
}