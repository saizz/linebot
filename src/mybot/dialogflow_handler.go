package mybot

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() {
	http.HandleFunc("/dialogflow", handleDialogflow)

	res := &response{
		Speech:      "error",
		DisplayText: "error",
	}

	errorResponse, _ = json.Marshal(res)
}

var (
	errorResponse []byte
)

type dialogflowRequest struct {
	ID     string `json:"id"`
	Result result `json:"result"`
}

type result struct {
	ResolvedQuery    string            `json:"resolvedQuery"`
	Action           string            `json:"action"`
	ActionIncomplete bool              `json:"actionIncomplete"`
	Parameters       map[string]string `json:"parameters"`
	Metadata         metadata          `json:"metadata"`
}

type metadata struct {
	IntentID   string `json:"intentId"`
	IntentName string `json:"intentName"`
}

type response struct {
	Speech      string `json:"speech"`
	DisplayText string `json:"displayText"`
}

// handleDialogflow is handle webhook from dialogflow.
func handleDialogflow(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_, _ = w.Write(toResponse(err.Error(), err.Error()))
		return
	}
	defer r.Body.Close()

	log.Infof(ctx, "dialogflow POST body: %v", string(b))

	req := new(dialogflowRequest)
	if err := json.Unmarshal(b, req); err != nil {
		_, _ = w.Write(toResponse(err.Error(), err.Error()))
		return
	}

	calID, err := resolveCalendar(req)
	if err != nil {
		_, _ = w.Write(toResponse(err.Error(), err.Error()))
		return
	}

	events, err := getCalendarEventsInternal(ctx, calID)
	if err != nil {
		_, _ = w.Write(toResponse(err.Error(), err.Error()))
		return
	}

	var text string
	for _, e := range events {
		text += e.Start.Date + e.Start.DateTime + "は" + e.Summary + "。\n"
	}

	_, _ = w.Write(toResponse(text, text))
}

func resolveCalendar(i *dialogflowRequest) (string, error) {

	if cal, ok := i.Result.Parameters["calendar"]; ok {
		if id := getCalendarID(cal); id != "" {
			return id, nil
		}
	}

	return "", errors.New("どのカレンダーがわかりませんでした。")
}

func getCalendarID(cal string) string {

	switch cal {
	case "holiday":
		return jaHolidayCal
	case "fg-space":
		return fgSpaceCal
	case "chigasaki-login":
		return chigaLoginCal
	}

	return ""
}

func toResponse(speech, text string) []byte {

	res := &response{
		Speech:      speech,
		DisplayText: text,
	}

	b, err := json.Marshal(res)
	if err != nil {
		return errorResponse
	}

	return b
}
