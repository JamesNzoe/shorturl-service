package services

import (
	hashids "github.com/speps/go-hashids"
)

const (
	defaultHashIDAlphabet  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	defaultHashIDSalt      = ",ebjyChTO021d7j0sVTLUQ90u5R0;"
	defaultHashIDMinLength = 2
)

func NewHasher() (*hashids.HashID, error) {
	hd := hashids.NewData()
	hd.Alphabet = defaultHashIDAlphabet
	hd.Salt = defaultHashIDSalt
	hd.MinLength = defaultHashIDMinLength

	d, err := hashids.NewWithData(hd)
	if err != nil {
		return nil, err
	}
	return d, nil
}
