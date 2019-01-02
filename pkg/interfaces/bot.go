package interfaces

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

type Bot interface {
	Send(to tb.Recipient, what interface{}, options ...interface{}) (*tb.Message, error)
	Handle(endpoint interface{}, handler interface{})
}
