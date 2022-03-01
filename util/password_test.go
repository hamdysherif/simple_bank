package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateHashedPassowrd(t *testing.T) {
	password, err := GenerateHashedPassowrd("123456")
	require.NoError(t, err)
	require.NotEmpty(t, password)
}

func TestCheckHashedPassword(t *testing.T) {

	hashedPassword, _ := GenerateHashedPassowrd("123456")

	testCases := map[string]struct {
		password       string
		hashedPassword string
		result         bool
	}{
		"ValidPassword": {
			password:       "123456",
			hashedPassword: hashedPassword,
			result:         true,
		},
		"InvalidPassword": {
			password:       "12345",
			hashedPassword: hashedPassword,
			result:         false,
		},
	}

	for i, tc := range testCases {
		t.Run(i, func(t *testing.T) {
			require.Equal(t, tc.result, CheckHashedPassword(tc.hashedPassword, tc.password))
		})
	}
}
