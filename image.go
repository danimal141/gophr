package main

import (
	"image"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/disintegration/imaging"
)

const imageIDLength = 10

// A map of accepted mime types and their file extension
var mimeExtensions = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/gif":  ".gif",
}

type Image struct {
	ID          string
	UserID      string
	Name        string
	Location    string
	Size        int64
	CreatedAt   time.Time
	Description string
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func (img *Image) ShowRoute() string {
	return "/image/" + img.ID
}

func (img *Image) StaticRoute() string {
	return "/im/" + img.Location
}

func NewImage(user *User) *Image {
	return &Image{
		ID:        GenerateID("img", imageIDLength),
		UserID:    user.ID,
		CreatedAt: time.Now(),
	}
}

func (img *Image) CreateFromURL(imgURL string) error {
	resp, err := http.Get(imgURL)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errImageURLInvalid
	}
	defer resp.Body.Close()

	mimeType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return errInvalidImageType
	}

	ext, ok := mimeExtensions[mimeType]
	if !ok {
		return errInvalidImageType
	}

	img.Name = filepath.Base(imgURL)
	img.Location = img.ID + ext
	file, err := os.Create("./data/images/" + img.Location)
	if err != nil {
		return err
	}
	defer file.Close()

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	img.Size = size

	err = img.CreateResizedImages()
	if err != nil {
		return err
	}
	return globalImageStore.Save(img)
}

func (img *Image) CreateFromFile(file multipart.File, headers *multipart.FileHeader) error {
	img.Name = headers.Filename
	img.Location = img.ID + filepath.Ext(img.Name)
	savedFile, err := os.Create("./data/images/" + img.Location)
	if err != nil {
		return err
	}
	defer file.Close()

	size, err := io.Copy(savedFile, file)
	if err != nil {
		return err
	}
	img.Size = size

	err = img.CreateResizedImages()
	if err != nil {
		return err
	}
	return globalImageStore.Save(img)
}

func (img *Image) CreateResizedImages() error {
	srcImg, err := imaging.Open("./data/images/" + img.Location)
	if err != nil {
		return err
	}

	errChan := make(chan error)
	go img.resizePreview(errChan, srcImg)
	go img.resizeThumbnail(errChan, srcImg)

	// Wait for images to finish resizing
	for i := 0; i < 2; i++ {
		err := <-errChan
		if err != nil {
			return err
		}
	}
	return nil
}

var widthPreview = 800
var widthThumbnail = 400

func (img *Image) resizePreview(errChan chan error, srcImg image.Image) {
	size := srcImg.Bounds().Size()
	ratio := float64(size.Y) / float64(size.X)
	height := int(float64(widthPreview) * ratio)

	dstImg := imaging.Resize(srcImg, widthPreview, height, imaging.Lanczos)
	dest := "./data/images/preview/" + img.Location
	errChan <- imaging.Save(dstImg, dest)
}

func (img *Image) resizeThumbnail(errChan chan error, srcImg image.Image) {
	dstImg := imaging.Thumbnail(srcImg, widthThumbnail, widthThumbnail, imaging.Lanczos)
	dest := "./data/images/thumbnail/" + img.Location
	errChan <- imaging.Save(dstImg, dest)
}
