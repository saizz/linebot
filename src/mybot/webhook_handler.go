package mybot

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
	"google.golang.org/appengine/urlfetch"
)

var webhookHandler *httphandler.WebhookHandler

func init() {
	err := godotenv.Load("line.env")
	if err != nil {
		panic(err)
	}

	webhookHandler, err = httphandler.New(
		os.Getenv("LINE_BOT_CHANNEL_SECRET"),
		os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
	)
	if err != nil {
		panic(err)
	}

	webhookHandler.HandleEvents(handleWebhook)
	http.Handle("/webhook", webhookHandler)
}

// newLineBot create linebot Client.
func newLineBot(c context.Context) (*linebot.Client, error) {
	return webhookHandler.NewClient(
		linebot.WithHTTPClient(urlfetch.Client(c)),
	)
}

// handleCallback handle webhook from linebot server and regist task queue.
func handleWebhook(events []*linebot.Event, r *http.Request) {
	ctx := appengine.NewContext(r)
	tasks := make([]*taskqueue.Task, 0, len(events))

	for _, e := range events {
		t, err := newPOSTJSONTask("/task", e)
		if err != nil {
			continue
		}
		tasks = append(tasks, t)
	}

	if _, err := taskqueue.AddMulti(ctx, tasks, "linebot-worker"); err != nil {
		log.Errorf(ctx, "taskqueue.AddMulti: %v", err)
	}
}

func newPOSTJSONTask(path string, i interface{}) (*taskqueue.Task, error) {
	h := make(http.Header)
	h.Set("Content-Type", "application/json")

	j, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	return &taskqueue.Task{
		Path:    path,
		Payload: j,
		Header:  h,
		Method:  "POST",
	}, nil
}
