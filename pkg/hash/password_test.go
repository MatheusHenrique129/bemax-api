package hash_test

import (
	"fmt"
	"testing"

	"github.com/MatheusHenrique129/bemax-api/pkg/hash"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	testCases := []struct {
		name          string
		password      string
		expectedError string
		isExpectedErr bool
	}{
		{
			name:          "success",
			password:      "password123",
			isExpectedErr: false,
		},
		{
			name:          "error",
			password:      "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$%^&*()-_=+{}[]",
			expectedError: fmt.Sprintf("error hashing password: %v", bcrypt.ErrPasswordTooLong),
			isExpectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := hash.HashPassword(tc.password)

			if tc.isExpectedErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestCheckPassword(t *testing.T) {
	testCases := []struct {
		name          string
		password      string
		passwordHash  string
		expectedError bool
	}{
		{
			name:          "success",
			password:      "password123",
			passwordHash:  "$2a$10$NI6tE59sN6VuTfNkbFOXFOAHr.4Oso/UtiL0tGjr2udPkFbM9AOU6",
			expectedError: false,
		},
		{
			name:          "error",
			password:      "passwordError",
			passwordHash:  "$2a$10$NI6tE59sN6VuTfNkbFOXFOAHr.4Oso/UtiL0tGjr2udPkFbM9AOU6",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := hash.CheckPassword(tc.password, tc.passwordHash)

			if tc.expectedError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
