package mybot

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

func init() {
	http.HandleFunc("/task", handleTask)
}

// handleTask hadle task.
func handleTask(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(r)
	data := r.FormValue("data")
	if data == "" {
		errorf(ctx, "No data")
		return
	}

	j, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		errorf(ctx, "base64 DecodeString: %v", err)
		return
	}

	e := new(linebot.Event)
	err = json.Unmarshal(j, e)
	if err != nil {
		errorf(ctx, "json.Unmarshal: %v", err)
		return
	}
	logf(ctx, "EventType: %s\nMessage: %#v", e.Type, e.Message)

	m := []linebot.Message{linebot.NewTextMessage("ok")}
	// wker := NewWorker(ctx, e)
	// if wker == nil {
	// 	return
	// }
	// m := wker.Reply()

	bot, err := newLineBot(ctx)
	if err != nil {
		errorf(ctx, "newLineBot: %v", err)
		return
	}

	if _, err = bot.ReplyMessage(e.ReplyToken, m...).WithContext(ctx).Do(); err != nil {
		errorf(ctx, "ReplayMessage: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
