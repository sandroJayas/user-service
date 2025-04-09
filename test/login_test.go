package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const baseURL = "http://localhost:8080"

func TestLogin(t *testing.T) {
	timestamp := time.Now().Format("150405")
	email := "login+" + timestamp + "@test.com"
	password := "strongpassword"

	// First register a user
	t.Run("setup - register user", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"email":    email,
			"password": password,
		})
		resp, err := http.Post(baseURL+"/users/register", "application/json", bytes.NewReader(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	// Now test login cases
	tests := []struct {
		name           string
		body           map[string]string
		expectedStatus int
		expectToken    bool
	}{
		{
			name: "valid login",
			body: map[string]string{
				"email":    email,
				"password": password,
			},
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
		{
			name: "wrong password",
			body: map[string]string{
				"email":    email,
				"password": "wrongpass",
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "non-existent user",
			body: map[string]string{
				"email":    "ghost+" + timestamp + "@test.com",
				"password": "doesntmatter",
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "missing email",
			body: map[string]string{
				"password": password,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing password",
			body: map[string]string{
				"email": email,
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		tc := tc // capture
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			resp, err := http.Post(baseURL+"/users/login", "application/json", bytes.NewReader(body))

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectToken {
				var res map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&res)
				assert.NoError(t, err)
				token := res["token"]
				assert.NotEmpty(t, token, "token should be present")
				assert.IsType(t, "", token, "token should be a string")
			}
		})
	}
}
