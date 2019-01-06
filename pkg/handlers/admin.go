package handlers

import (
	"github.com/vincenciusgeraldo/sibyl"
	tb "gopkg.in/tucnak/telebot.v2"
	"github.com/vincenciusgeraldo/sibyl/pkg/utils"
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
