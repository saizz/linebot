package mybot

import (
	"golang.org/x/net/context"

	"github.com/line/line-bot-sdk-go/linebot"
)

// Worker has Reply method which return linebot.Message interface.
type Worker interface {
	Reply() []linebot.Message
}

// NewWorker create new Worker.
func NewWorker(c context.Context, e *linebot.Event) Worker {

	switch m := e.Message.(type) {
	case *linebot.TextMessage:
		return NewTextWorker(m)
	case *linebot.ImageMessage:
		return NewImageWorker(c, m)
	}

	return nil
}
