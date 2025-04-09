package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateProfile(t *testing.T) {
	email := "profile+" + time.Now().Format("150405") + "@test.com"
	password := "supersecure"

	var token string

	t.Run("register user", func(t *testing.T) {
		payload := map[string]string{
			"email":    email,
			"password": password,
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

	t.Run("update profile", func(t *testing.T) {
		update := map[string]string{
			"first_name":        "New",
			"last_name":         "Name",
			"address_line_1":    "456 New Ave",
			"address_line_2":    "Suite 100",
			"city":              "Newtown",
			"postal_code":       "12345",
			"country":           "Testland",
			"phone_number":      "9876543210",
			"payment_method_id": "pm_123456",
		}
		body, _ := json.Marshal(update)

		req, _ := http.NewRequest("PUT", baseURL+"/users/profile", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		user := res["user"].(map[string]interface{})

		for key, _ := range update {
			assert.Equal(t, update[key], user[key])
		}
	})

}
