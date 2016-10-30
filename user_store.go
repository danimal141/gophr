package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

type UserStore interface {
	Find(string) (*User, error)
	FindByEmail(string) (*User, error)
	FindByUsername(string) (*User, error)
	Save(User) error
}

type FileUserStore struct {
	filename string
	Users    map[string]User
}

func (store FileUserStore) Find(id string) (*User, error) {
	user, ok := store.Users[id]
	if ok {
		return &user, nil
	}
	return nil, nil
}

func (store FileUserStore) FindByEmail(email string) (*User, error) {
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

func (store FileUserStore) FindByUsername(username string) (*User, error) {
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

func (store FileUserStore) Save(user User) error {
	store.Users[user.ID] = user

	contents, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(store.filename, contents, 0660)
	if err != nil {
		return err
	}
	return nil
}
