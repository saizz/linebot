package mybot

import "github.com/line/line-bot-sdk-go/linebot"

// TextWorker is Worker for TextMessage.
type TextWorker struct {
	message *linebot.TextMessage
}

// NewTextWorker create new Worker.
func NewTextWorker(m *linebot.TextMessage) Worker {
	return &TextWorker{
		message: m,
	}
}

// Reply return linebot.Message interface.
func (w *TextWorker) Reply() []linebot.Message {
	m := make([]linebot.Message, 0, 1)
	m = append(m, linebot.NewTextMessage(w.message.Text+"???"))
	return m
}
