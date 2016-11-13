package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

type UserStore interface {
	Find(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByUsername(username string) (*User, error)
	Save(user *User) error
}

type FileUserStore struct {
	filename string
	Users    map[string]User
}

var globalUserStore UserStore

func NewFileUserStore(filename string) (*FileUserStore, error) {
	store := &FileUserStore{filename: filename, Users: map[string]User{}}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return store, err
		}
		return nil, err
	}

	err = json.Unmarshal(contents, store)
	if err != nil {
		return nil, err
	}
	return store, nil
}

func (store *FileUserStore) Find(id string) (*User, error) {
	user, ok := store.Users[id]
	if ok {
		return &user, nil
	}
	return nil, nil
}

func (store *FileUserStore) FindByEmail(email string) (*User, error) {
	if email == "" {
		return nil, nil
	}
	for _, user := range store.Users {
		if strings.ToLower(user.Email) == strings.ToLower(email) {
			return &user, nil
		}
	}
	return nil, nil
}

func (store *FileUserStore) FindByUsername(username string) (*User, error) {
	if username == "" {
		return nil, nil
	}
	for _, user := range store.Users {
		if strings.ToLower(user.Username) == strings.ToLower(username) {
			return &user, nil
		}
	}
	return nil, nil
}

func (store *FileUserStore) Save(user *User) error {
	store.Users[user.ID] = *user

	contents, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(store.filename, contents, 0660)
}
