package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type Pasetomaker struct {
	symmetricKey []byte
	paseto       *paseto.V2
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) < chacha20poly1305.KeySize {
		return nil, fmt.Errorf("secret symmetricKey length should not be less than %d", chacha20poly1305.KeySize)
	}
	return &Pasetomaker{
		symmetricKey: []byte(symmetricKey),
		paseto:       paseto.NewV2(),
	}, nil
}

// CreateToken create and sign a token for username with duration
func (maker *Pasetomaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPaylod(username, duration)
	if err != nil {
		return "", err
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		return "", err
	}
	return token, nil
}

// VerifyToken verify the token and return the decoded payload
func (maker *Pasetomaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
