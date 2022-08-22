package utils

import (
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

// Wraps copier.Copy(...) (github.com/jinzhu/copier)
//
// Just because copy from source to destination feels more natural
func Copy(from interface{}, to interface{}) error {
	return copier.Copy(to, from)
}

// Wraps copier.CopyWithOption(...) (github.com/jinzhu/copier)
//
// Just because copy from source to destination feels more natural
func CopyWithOption(from interface{}, to interface{}, opt copier.Option) error {
	return copier.CopyWithOption(to, from, opt)
}

// Create a password hash
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	return string(bytes), err
}

// Check a hash against a password
func VerifyPassword(hash string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err == nil {
		return true
	}
	return false
}
