package mybot

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/context"

	"google.golang.org/appengine/taskqueue"
	"google.golang.org/appengine/urlfetch"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
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
	webhookHandler.HandleEvents(handleWebhook)
	http.Handle("/webhook", webhookHandler)
}

// handleCallback handle webhook from linebot server and regist task queue.
func handleWebhook(events []*linebot.Event, r *http.Request) {
	ctx := newContext(r)
	tasks := make([]*taskqueue.Task, len(events))

	for i, e := range events {
		j, err := json.Marshal(e)
		if err != nil {
			errorf(ctx, "json.Marshal: %v", err)
			return
		}
		data := base64.StdEncoding.EncodeToString(j)
		t := taskqueue.NewPOSTTask("/task", url.Values{"data": {data}})
		tasks[i] = t
	}

	if _, err := taskqueue.AddMulti(ctx, tasks, "default"); err != nil {
		errorf(ctx, "taskqueue.AddMulti: %v", err)
	}
}

// newLineBot create linebot Client.
func newLineBot(c context.Context) (*linebot.Client, error) {
	return webhookHandler.NewClient(
		linebot.WithHTTPClient(urlfetch.Client(c)),
	)
}
