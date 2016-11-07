package main

type ImageStore interface {
	Save(img *Image) error
	Find(id string) (*Image, error)
	FindAll(offset int) ([]Image, error)
	FindAllByUser(user *User, offset int) ([]Image, error)
}
