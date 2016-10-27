package main

import (
	"crypto/rand"
	"fmt"
)

const idSource = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const idSourceLen = byte(len(idSource))

func GenerateID(prefix string, length int) (string, error) {
	id := make([]byte, length)

	_, err := rand.Read(id)
	if err != nil {
		return "", err
	}

	// Replace each random number with an alphanumeric value
	for i, b := range id {
		id[i] = idSource[b%idSourceLen]
	}

	return fmt.Sprintf("%s_%s", prefix, string(id)), nil
}
