package mybot

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"

	"cloud.google.com/go/storage"
	"github.com/line/line-bot-sdk-go/linebot"
	"google.golang.org/appengine"
	"google.golang.org/appengine/blobstore"
	"google.golang.org/appengine/file"
	gaeimage "google.golang.org/appengine/image"
	"google.golang.org/appengine/log"
)

// ImageWorker is Worker for ImageMessage.
type ImageWorker struct {
	message *linebot.ImageMessage
}

// NewImageWorker create new Worker.
func NewImageWorker(m *linebot.ImageMessage) Worker {
	return &ImageWorker{
		message: m,
	}
}

// Reply return linebot.Message interface.
func (w *ImageWorker) Reply(ctx context.Context) []linebot.Message {
	m := make([]linebot.Message, 0, 2)

	img, err := w.getImageContent(ctx)
	if err != nil {
		log.Errorf(ctx, "getImageContent: %v", err)
		m = append(m, linebot.NewTextMessage("cant get image."))
		return m
	}

	err = w.storeImage(ctx, toGrayscale(img))
	if err != nil {
		log.Errorf(ctx, "storeImage: %v", err)
		m = append(m, linebot.NewTextMessage("cant save storeage."))
		return m
	}

	m = append(m, linebot.NewImageMessage(
		w.getConvertedImageURL(ctx),
		w.getConvertedPreviewImageURL(ctx)))

	servingURL, err := w.getServingURL(ctx)
	if err != nil {
		m = append(m, linebot.NewTextMessage("cant get serving URL."))
	}

	m = append(m, linebot.NewImageMessage(
		servingURL,
		servingURL+"=s128-cc"))

	m = append(m, linebot.NewTextMessage("covert done."))
	return m
}

// getImageContent get Message Content from linebot server.
func (w *ImageWorker) getImageContent(ctx context.Context) (image.Image, error) {

	bot, err := newLineBot(ctx)
	if err != nil {
		return nil, err
	}

	res, err := bot.GetMessageContent(w.message.ID).WithContext(ctx).Do()
	if err != nil {
		return nil, err
	}
	defer res.Content.Close()

	img, err := jpeg.Decode(res.Content)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// storeImage store imge to GCS.
func (w *ImageWorker) storeImage(ctx context.Context, img image.Image) error {

	bucket, err := file.DefaultBucketName(ctx)
	if err != nil {
		return err
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	writer := client.Bucket(bucket).Object(w.getObjectName()).NewWriter(ctx)
	writer.ContentType = "image/jpeg"

	if err := jpeg.Encode(writer, img, nil); err != nil {
		return err
	}

	return writer.Close()
}

// getObjectName return GCS object name.
func (w *ImageWorker) getObjectName() string {
	return w.message.ID + ".jpeg"
}

func (w *ImageWorker) getServingURL(ctx context.Context) (string, error) {

	bucket, err := file.DefaultBucketName(ctx)
	if err != nil {
		return "", err
	}

	gsURL := fmt.Sprintf("/gs/%s/%s", bucket, w.getObjectName())
	blobKey, err := blobstore.BlobKeyForFile(ctx, gsURL)
	if err != nil {
		return "", err
	}

	opts := &gaeimage.ServingURLOptions{Secure: true}
	url, err := gaeimage.ServingURL(ctx, blobKey, opts)
	if err != nil {
		return "", err
	}
	log.Infof(ctx, "servingURL: %v", url)
	return url.String(), nil
}

// getConvertedImageUrl return converted image url.
func (w *ImageWorker) getConvertedImageURL(ctx context.Context) string {
	id := appengine.AppID(ctx)

	return "https://" + id + ".appspot.com/image?name=" + w.message.ID + ".jpeg"
}

// getConvertedPreviewImageURL retrun converted preview image url.
func (w *ImageWorker) getConvertedPreviewImageURL(ctx context.Context) string {
	return w.getConvertedImageURL(ctx) + "&preview=true"
}

// Converted implements image.Image, so you can
// pretend that it is the converted image.
type Converted struct {
	Img image.Image
	Mod color.Model
}

// ColorModel return image's color.Model.
func (c *Converted) ColorModel() color.Model {
	return c.Mod
}

// Bounds retrun image's bounds.
func (c *Converted) Bounds() image.Rectangle {
	return c.Img.Bounds()
}

// At forwards the call to the original image and
// then asks the color model to convert it.
func (c *Converted) At(x, y int) color.Color {
	return c.Mod.Convert(c.Img.At(x, y))
}

// toGrayscale retrun image.Image interface converted grayscale.
func toGrayscale(src image.Image) image.Image {
	return &Converted{src, color.GrayModel}
}
