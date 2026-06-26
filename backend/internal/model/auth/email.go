package model

import (
	"crypto/rand"
	"encoding/base64"
)

type Message struct {
	From    string
	To      []string
	Html    string
	Subject string
}

type TokenVerify struct {
	TokenVer string
}

func NewTokenVerify() (TokenVerify, error) {
	temp := make([]byte, 32)

	_, err := rand.Read(temp)
	if err != nil {
		return TokenVerify{}, err
	}

	return TokenVerify{TokenVer: base64.URLEncoding.EncodeToString(temp)}, nil
}
