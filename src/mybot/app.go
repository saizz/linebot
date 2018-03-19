package mybot

import (
	"context"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// newContext is wrapper of appengine.NewContext.
func newContext(r *http.Request) context.Context {
	return appengine.NewContext(r)
}

// logf is wrapper of log.Infof.
func logf(c context.Context, format string, args ...interface{}) {
	log.Infof(c, format, args...)
}

// errorf is wrapper of log.Errorf.
func errorf(c context.Context, format string, args ...interface{}) {
	log.Errorf(c, format, args...)
}
