package password

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/giantswarm/microerror"
	"golang.org/x/crypto/bcrypt"
)

type Manager interface {
	GeneratePassword(length int) (string, error)
	Hash(password []byte) ([]byte, error)
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

func (m SimpleManager) Hash(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(password, 14)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return hash, nil
}
