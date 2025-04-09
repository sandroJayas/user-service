package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("DATABASE_URL", "postgres://testuser:testpass@localhost:5433/user_service_test?sslmode=disable")
	os.Setenv("JWT_SECRET", "test-secret")

	time.Sleep(2 * time.Second) // or wait for /healthz

	os.Exit(m.Run())
}

func TestRegister(t *testing.T) {
	timestamp := time.Now().Format("150405")

	tests := []struct {
		name           string
		body           map[string]string
		expectedStatus int
		expectUser     bool
	}{
		{
			name: "valid user",
			body: map[string]string{
				"email":    "valid+" + timestamp + "@test.com",
				"password": "supersecret",
			},
			expectedStatus: http.StatusCreated,
			expectUser:     true,
		},
		{
			name: "missing email",
			body: map[string]string{
				"password": "test123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid email format",
			body: map[string]string{
				"email":    "not-an-email",
				"password": "test123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing password",
			body: map[string]string{
				"email": "no-pass+" + timestamp + "@test.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "duplicate email",
			body: map[string]string{
				"email":    "valid+" + timestamp + "@test.com",
				"password": "test123",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)

			resp, err := http.Post(baseURL+"/users/register", "application/json", bytes.NewReader(body))
			assert.NoError(t, err, "request should not error")
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectUser {
				var res map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&res)
				assert.NoError(t, err)
				user, ok := res["user"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, tc.body["email"], user["email"])
			}
		})
	}
}
