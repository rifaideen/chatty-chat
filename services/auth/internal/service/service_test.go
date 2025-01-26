package service

import (
	"pkg/auth"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	service := New("secret", 1)

	t.Run("auth error", func(t *testing.T) {
		token, err := service.Login(&auth.LoginRequest{
			Username: "admin",
			Password: "",
		})

		assert.Empty(t, token)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "invalid credentials")
	})

	t.Run("auth success", func(t *testing.T) {
		token, err := service.Login(&auth.LoginRequest{
			Username: "admin",
			Password: "admin",
		})

		assert.NotEmpty(t, token)
		assert.NoError(t, err)
	})
}

func TestVerifyToken(t *testing.T) {
	service := New("secret", 1)

	t.Run("verify token error", func(t *testing.T) {
		token, err := service.VerifyToken(&auth.VerifyRequest{
			Token: "invalid token",
		})

		assert.Nil(t, token)
		assert.Error(t, err)
	})

	t.Run("verify token success", func(t *testing.T) {
		token, err := service.Login(&auth.LoginRequest{
			Username: "admin",
			Password: "admin",
		})

		assert.NotEmpty(t, token)
		assert.NoError(t, err)

		verifiedToken, err := service.VerifyToken(&auth.VerifyRequest{
			Token: token,
		})

		assert.NotNil(t, verifiedToken)
		assert.NoError(t, err)
	})
}
