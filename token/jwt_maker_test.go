package token

import (
	"testing"
	"time"

	"github.com/hamdysherif/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestGenerateJWTToken(t *testing.T) {
	maker, err := NewJWTMaker("12345678901234567890123456789023")
	require.NoError(t, err)
	username := util.RandomOwner()

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(time.Minute)
	token, err := maker.CreateToken(username, time.Minute)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.Equal(t, payload.Username, username)
	require.NotZero(t, payload.ID)
	require.WithinDuration(t, payload.ExpireAt, expiredAt, time.Second)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
}

func TestExpireJWTToken(t *testing.T) {
	maker, err := NewJWTMaker("12345678901234567890123456789023")
	require.NoError(t, err)
	username := util.RandomOwner()

	token, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NotEmpty(t, err)
	require.Nil(t, payload)
	require.Error(t, err, ErrExpiredToken)
}
