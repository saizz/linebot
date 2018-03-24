package mybot

import (
	"context"

	"github.com/line/line-bot-sdk-go/linebot"
)

// Worker has Reply method which return linebot.Message interface.
type Worker interface {
	Reply(context.Context) []linebot.Message
}

// NewWorker create new Worker.
func NewWorker(e *linebot.Event) Worker {

	switch m := e.Message.(type) {
	case *linebot.TextMessage:
		return NewTextWorker(m)
	case *linebot.ImageMessage:
		return NewImageWorker(m)
	}

	return nil
}
