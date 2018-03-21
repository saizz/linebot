package mybot

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

func init() {
	http.HandleFunc("/task", handleTask)
}

// handleTask hadle task.
func handleTask(w http.ResponseWriter, r *http.Request) {

	ctx := newContext(r)
	e, err := getEvent(r)
	if err != nil {
		errorf(ctx, "getEvent: %v", err)
		return
	}

	logf(ctx, "EventType: %s\nMessage: %#v", e.Type, e.Message)

	//m := []linebot.Message{linebot.NewTextMessage("ok")}
	wk := NewWorker(ctx, e)
	if wk == nil {
		return
	}
	m := wk.Reply()

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

// getEvent get linebot Event from http request.
func getEvent(r *http.Request) (e *linebot.Event, err error) {

	var b []byte
	if b, err = ioutil.ReadAll(r.Body); err != nil {
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(b, e)

	return
}
