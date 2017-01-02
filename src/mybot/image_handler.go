package mybot

import (
	"io/ioutil"
	"net/http"

	"cloud.google.com/go/storage"

	"google.golang.org/appengine"
	"google.golang.org/appengine/file"
)

func init() {
	http.HandleFunc("/image", handleImage)
}

// handleImage hadle image.
func handleImage(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")

	ctx := appengine.NewContext(r)
	bucket, err := file.DefaultBucketName(ctx)
	if err != nil {
		errorf(ctx, "handleImage DefaultBucketName: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		errorf(ctx, "handleImage storage.NewClient: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer client.Close()

	reader, err := client.Bucket(bucket).Object(name).NewReader(ctx)
	if err != nil {
		errorf(ctx, "handleImage Bucket NewReader: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		errorf(ctx, "handleImage Bucket ReadAll: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		errorf(ctx, "handleImage ResponseWriter Write: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
}
