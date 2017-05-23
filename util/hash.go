package util

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

// MD5sum MD5 hash
func MD5sum(content string) string {
	h := md5.New()
	io.WriteString(h, content)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SHA256 hash
func SHA256(content string) string {
	h := sha256.New()
	io.WriteString(h, content)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// BCrypt bcrypt hash
// https://en.wikipedia.org/wiki/Bcrypt
// https://godoc.org/golang.org/x/crypto/bcrypt
func BCrypt(password string) (b string, err error) {
	data, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(data), err
}

// CompareBCryptPassword compares a bcrypt hashed password with its possible plaintext
func CompareBCryptPassword(hashPassword, password string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}
