package entities

import (
	"errors"
	"regexp"
)

type Client string

var clientNameCheck = regexp.MustCompile("^[a-z0-9-_]+$").MatchString

func NewClientName(raw string) (Client, error) {
	if clientNameCheck(raw) {
		return Client(raw), nil
	}
	return "", errors.New("invalid client name")
}
