package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProtectedRouteAccess(t *testing.T) {
	email := "me+" + time.Now().Format("150405") + "@test.com"
	password := "testpass"

	t.Run("setup - register", func(t *testing.T) {
		payload := map[string]string{
			"email":    email,
			"password": password,
		}
		body, _ := json.Marshal(payload)
		resp, err := http.Post(baseURL+"/users/register", "application/json", bytes.NewReader(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	var token string

	t.Run("setup - login", func(t *testing.T) {
		payload := map[string]string{
			"email":    email,
			"password": password,
		}
		body, _ := json.Marshal(payload)
		resp, err := http.Post(baseURL+"/users/login", "application/json", bytes.NewReader(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]string
		_ = json.NewDecoder(resp.Body).Decode(&res)
		token = res["token"]
		assert.NotEmpty(t, token)
	})

	t.Run("valid token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/users/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := http.DefaultClient.Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			User struct {
				ID string `json:"id"`
			} `json:"user"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&body)
		assert.Contains(t, body.User.ID, "-", "should return a UUID user_id")
	})

	t.Run("missing token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/users/me", nil)
		resp, err := http.DefaultClient.Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("invalid token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/users/me", nil)
		req.Header.Set("Authorization", "Bearer not.a.valid.token")
		resp, err := http.DefaultClient.Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var body map[string]string
		_ = json.NewDecoder(resp.Body).Decode(&body)
		assert.True(t, strings.Contains(body["error"], "Invalid token"))
	})
}
