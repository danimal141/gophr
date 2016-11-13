package main

import (
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
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
	return globalImageStore.Save(img)
}

func (img *Image) ShowRoute() string {
	return "/image/" + img.ID
}

func (img *Image) StaticRoute() string {
	return "/im/" + img.Location
}
