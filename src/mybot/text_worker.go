package mybot

import (
	"context"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
	"google.golang.org/appengine/log"
)

const (
	// jaHolidayCal is japanese holiday calendar
	jaHolidayCal = "japanese__ja@holiday.calendar.google.com"
	// fgSpaceCal is FG-Space calendar
	fgSpaceCal = "freegufo.com_63de9ce3hvo4eruplc0sg8ii50@group.calendar.google.com"
	// chigaLoginCal is chigasaki LOGIN calendar
	chigaLoginCal = "wifi.de.login@gmail.com"
	// MonthSuffix is End time after 'MonthSuffix' months
	MonthSuffix = 6
	// MaxCalendarEventSize is max size getting from google calendar
	MaxCalendarEventSize = 5
	// TimeZone is Timezone at tokyo
	TimeZone = "Asia/Tokyo"
	// TimeZoneSuffix is Timezone suffix
	TimeZoneSuffix = 9 * 60 * 60
)

// TextWorker is Worker for TextMessage.
type TextWorker struct {
	ctx     context.Context
	message *linebot.TextMessage
}

// NewTextWorker create new Worker.
func NewTextWorker(c context.Context, m *linebot.TextMessage) Worker {
	return &TextWorker{
		ctx:     c,
		message: m,
	}
}

// Reply return linebot.Message interface.
func (w *TextWorker) Reply() []linebot.Message {
	m := make([]linebot.Message, 0, MaxCalendarEventSize)
	if !strings.Contains(w.message.Text, "祝日") {
		m = append(m, linebot.NewTextMessage(w.message.Text+"???"))
		return m
	}

	events, err := w.getCalendarEvents()
	if err != nil {
		log.Errorf(w.ctx, "getCalendarEvents: %v", err)
		m = append(m, linebot.NewTextMessage("cant get calendar."))
		return m
	}

	if len(events) == 0 {
		m = append(m, linebot.NewTextMessage("no 休み"))
		return m
	}

	for _, item := range events {
		m = append(m, linebot.NewTextMessage(item.Start.Date+"は"+item.Summary))
	}

	return m
}

// getCalendar get google calendar event.
func (w *TextWorker) getCalendarEvents() ([]*calendar.Event, error) {
	return getCalendarEventsInternal(w.ctx, jaHolidayCal)
}

func getCalendarEventsInternal(ctx context.Context, id string) ([]*calendar.Event, error) {

	client, err := google.DefaultClient(ctx, calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, err
	}

	svc, err := calendar.New(client)
	if err != nil {
		return nil, err
	}

	start := nowJST()
	end := start.AddDate(0, MonthSuffix, 0)
	log.Infof(ctx, "getCalendar start: %v, end: %v", start, end)

	events, err := svc.Events.List(id).
		TimeMin(start.Format(time.RFC3339)).
		TimeMax(end.Format(time.RFC3339)).
		TimeZone(TimeZone).
		SingleEvents(true).
		OrderBy("startTime").
		MaxResults(MaxCalendarEventSize).Do()
	if err != nil {
		return nil, err
	}

	return events.Items, nil

}

func nowJST() time.Time {
	jst := time.FixedZone(TimeZone, TimeZoneSuffix)
	return time.Now().UTC().In(jst)
}
