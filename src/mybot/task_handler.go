package mybot

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() {
	http.HandleFunc("/task", handleTask)
}

// handleTask hadle task.
func handleTask(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	event, err := getEvent(ctx, r)
	if err != nil {
		log.Errorf(ctx, "getEvent: %v", err)
		return
	}

	//msgs := []linebot.Message{linebot.NewTextMessage("ok")}
	worker := NewWorker(event)
	if worker == nil {
		return
	}
	msgs := worker.Reply(ctx)

	bot, err := newLineBot(ctx)
	if err != nil {
		log.Errorf(ctx, "newLineBot: %v", err)
		return
	}

	if _, err = bot.ReplyMessage(event.ReplyToken, msgs...).WithContext(ctx).Do(); err != nil {
		log.Errorf(ctx, "ReplayMessage: %v", err)
		return
	}
}

// getEvent get linebot Event from http request.
func getEvent(ctx context.Context, r *http.Request) (e *linebot.Event, err error) {

	var b []byte
	e = new(linebot.Event)
	if b, err = ioutil.ReadAll(r.Body); err != nil {
		return
	}
	defer r.Body.Close()

	log.Infof(ctx, "event: %v", string(b))
	err = json.Unmarshal(b, e)

	return
}
