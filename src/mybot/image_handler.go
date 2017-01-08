package mybot

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/nfnt/resize"

	"cloud.google.com/go/storage"

	"bytes"
	"image/jpeg"

	"google.golang.org/appengine"
	"google.golang.org/appengine/file"
)

const (
	// PreviewWidth is preview image width
	PreviewWidth = 160
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
	preview := r.FormValue("preview")

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
	defer reader.Close()

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		errorf(ctx, "handleImage Bucket ReadAll: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if preview == "true" {
		b = toPreview(b)
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	_, err = w.Write(b)
	if err != nil {
		errorf(ctx, "handleImage ResponseWriter Write: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// toPreview convert to preview image.
func toPreview(org []byte) []byte {

	img, err := jpeg.Decode(bytes.NewReader(org))
	if err != nil {
		return org
	}

	buf := new(bytes.Buffer)
	jpeg.Encode(buf,
		resize.Resize(PreviewWidth, 0, img, resize.Lanczos3), nil)

	return buf.Bytes()
}
