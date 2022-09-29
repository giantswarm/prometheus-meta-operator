package password

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/giantswarm/microerror"
	"golang.org/x/crypto/bcrypt"
)

type Manager interface {
	GeneratePassword(length int) (string, error)
	Hash(plaintext string) (string, error)
}

type SimpleManager struct {
}

func (m SimpleManager) GeneratePassword(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (m SimpleManager) Hash(plaintext string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), 14)
	if err != nil {
		return "", microerror.Mask(err)
	}
	return string(bytes), nil
}
