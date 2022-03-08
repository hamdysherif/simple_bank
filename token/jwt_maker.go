package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const secretKeyAllowedLength = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < secretKeyAllowedLength {
		return nil, fmt.Errorf("secret key length should not be less than %d", secretKeyAllowedLength)
	}
	return &JWTMaker{secretKey}, nil
}

func (jwtMaker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {

	payload, err := NewPaylod(username, duration)
	if err != nil {
		return "", err
	}
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(jwtMaker.secretKey))

	return tokenString, err
}

func (jwtMaker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	// Parse methods use this callback function to supply
	// the key for verification.  The function receives the parsed,
	// but unverified Token.  This allows you to use properties in the
	// Header of the token (such as `kid`) to identify which key to use.
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(jwtMaker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		return nil, err
	}

	return jwtToken.Claims.(*Payload), nil
}
