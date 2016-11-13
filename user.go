package main

import (
	"crypto/md5"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             string
	Username       string
	Email          string
	HashedPassword string
}

const (
	passwordMinLength = 8
	userIDLength      = 16
)

func NewUser(username, email, password string) (*User, error) {
	user := &User{Username: username, Email: email}

	if username == "" {
		return user, errNoUsername
	}
	if email == "" {
		return user, errNoEmail
	}
	if password == "" {
		return user, errNoPassword
	}
	if len(password) < passwordMinLength {
		return user, errPasswordTooShort
	}

	existingUser, err := globalUserStore.FindByUsername(username)
	if err != nil {
		return user, err
	}
	if existingUser != nil {
		return user, errUsernameExists
	}

	existingUser, err = globalUserStore.FindByEmail(email)
	if err != nil {
		return user, err
	}
	if existingUser != nil {
		return user, errEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	user.HashedPassword = string(hashedPassword)
	user.ID = GenerateID("usr", userIDLength)
	return user, err
}

func FindUser(username, password string) (*User, error) {
	u := &User{Username: username}
	existingUser, err := globalUserStore.FindByUsername(username)
	if err != nil {
		return u, err
	}
	if existingUser == nil {
		return u, errCredentialsIncorrect
	}
	if err = bcrypt.CompareHashAndPassword(
		[]byte(existingUser.HashedPassword),
		[]byte(password),
	); err != nil {
		return u, errCredentialsIncorrect
	}
	return existingUser, nil
}

func UpdateUser(user *User, email, currentPassword, newPassword string) (User, error) {
	u := *user // copy user
	u.Email = email

	existingUser, err := globalUserStore.FindByEmail(email)
	if err != nil {
		return u, err
	}
	if existingUser != nil && existingUser.ID != user.ID {
		return u, errEmailExists
	}

	user.Email = email

	if currentPassword == "" {
		return u, nil
	}

	if err = bcrypt.CompareHashAndPassword(
		[]byte(u.HashedPassword),
		[]byte(currentPassword),
	); err != nil {
		return u, errPasswordIncorrect
	}

	if newPassword == "" {
		return u, errNoPassword
	}

	if len(newPassword) < passwordMinLength {
		return u, errPasswordTooShort
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	user.HashedPassword = string(hashedPassword)
	return u, err
}

func (user *User) AvatarURL() string {
	return fmt.Sprintf(
		"//www.gravatar.com/avatar/%x",
		md5.Sum([]byte(user.Email)),
	)
}

func (user *User) ImagesRoute() string {
	return "/user/" + user.ID
}
