package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPassEncrpt(t *testing.T) {
	pwd := []byte("12345#hello")
	encrypted, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	assert.NoError(t, err)
	println(string(encrypted))

	err = bcrypt.CompareHashAndPassword(encrypted, []byte("123"))
	assert.NoError(t, err)
}
