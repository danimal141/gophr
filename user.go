package main

import "golang.org/x/crypto/bcrypt"

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

func NewUser(username, email, password string) (User, error) {
	user := User{Username: username, Email: email}

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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.HashedPassword = string(hashedPassword)
	user.ID = GenerateID("usr", userIDLength)
	return user, err
}
