package mybot

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"

	"cloud.google.com/go/storage"
	"github.com/line/line-bot-sdk-go/linebot"
	"google.golang.org/appengine"
	"google.golang.org/appengine/file"
)

// ImageWorker is Worker for ImageMessage.
type ImageWorker struct {
	ctx     context.Context
	message *linebot.ImageMessage
}

// NewImageWorker create new Worker.
func NewImageWorker(c context.Context, m *linebot.ImageMessage) Worker {
	return &ImageWorker{
		ctx:     c,
		message: m,
	}
}

// Reply return linebot.Message interface.
func (w *ImageWorker) Reply() []linebot.Message {
	m := make([]linebot.Message, 0, 2)

	img, err := w.getImageContent()
	if err != nil {
		errorf(w.ctx, "getImageContent: %v", err)
		m = append(m, linebot.NewTextMessage("cant get image."))
		return m
	}

	err = w.storeImage(toGrayscale(img))
	if err != nil {
		errorf(w.ctx, "storeImage: %v", err)
		m = append(m, linebot.NewTextMessage("cant save storeage."))
		return m
	}

	m = append(m, linebot.NewImageMessage(
		w.getConvertedImageURL(),
		w.getConvertedPreviewImageURL()))
	m = append(m, linebot.NewTextMessage("covert done."))
	return m
}

// getImageContent get Message Content from linebot server.
func (w *ImageWorker) getImageContent() (image.Image, error) {

	bot, err := newLineBot(w.ctx)
	if err != nil {
		return nil, err
	}

	res, err := bot.GetMessageContent(w.message.ID).WithContext(w.ctx).Do()
	if err != nil {
		return nil, err
	}
	defer res.Content.Close()

	b, err := ioutil.ReadAll(res.Content)
	if err != nil {
		return nil, err
	}

	img, err := jpeg.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	return img, nil
}

// storeImage store imge to GCS.
func (w *ImageWorker) storeImage(img image.Image) error {

	bucket, err := file.DefaultBucketName(w.ctx)
	if err != nil {
		return err
	}

	client, err := storage.NewClient(w.ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	writer := client.Bucket(bucket).Object(w.message.ID + ".jpeg").NewWriter(w.ctx)
	writer.ContentType = "image/jpeg"
	defer writer.Close()

	return jpeg.Encode(writer, img, nil)
}

// getConvertedImageUrl return converted image url.
func (w *ImageWorker) getConvertedImageURL() string {
	id := appengine.AppID(w.ctx)

	return "https://" + id + ".appspot.com/image?name=" + w.message.ID + ".jpeg"
}

// getConvertedPreviewImageURL retrun converted preview image url.
func (w *ImageWorker) getConvertedPreviewImageURL() string {
	return w.getConvertedImageURL() + "&preview=true"
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
