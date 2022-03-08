package token

import (
	"time"
)

type Maker interface {

	// CreateToken create and sign a token for username with duration
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken verify the token and return the decoded payload
	VerifyToken(token string) (*Payload, error)
}
