package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDeleteUser(t *testing.T) {
	email := "profile+" + time.Now().Format("150405") + "@test.com"
	password := "supersecure"

	var token string

	t.Run("register user", func(t *testing.T) {
		payload := map[string]string{
			"email":         email,
			"password":      password,
			"first_name":    "Old",
			"last_name":     "Name",
			"address_line1": "123 Old St",
			"city":          "Oldtown",
			"postal_code":   "00000",
			"country":       "Testland",
			"phone_number":  "1234567890",
		}
		body, _ := json.Marshal(payload)
		resp, err := http.Post(baseURL+"/users/register", "application/json", bytes.NewReader(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("login to get token", func(t *testing.T) {
		payload := map[string]string{
			"email":    email,
			"password": password,
		}
		body, _ := json.Marshal(payload)
		resp, err := http.Post(baseURL+"/users/login", "application/json", bytes.NewReader(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]string
		json.NewDecoder(resp.Body).Decode(&res)
		token = res["token"]
		assert.NotEmpty(t, token)
	})

	t.Run("delete user", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", baseURL+"/users/delete", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("verify access is blocked after delete", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/users/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("register same user", func(t *testing.T) {
		payload := map[string]string{
			"email":         email,
			"password":      password,
			"first_name":    "Old",
			"last_name":     "Name",
			"address_line1": "123 Old St",
			"city":          "Oldtown",
			"postal_code":   "00000",
			"country":       "Testland",
			"phone_number":  "1234567890",
		}
		body, _ := json.Marshal(payload)
		resp, err := http.Post(baseURL+"/users/register", "application/json", bytes.NewReader(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})
}
