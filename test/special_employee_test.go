package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSpecialEmployeeEndpoint(t *testing.T) {
	employeeEmail := "employee+" + time.Now().Format("150405") + "@sort.com"
	customerEmail := "customer+" + time.Now().Format("150405") + "@sort.com"
	password := "SuperSecure123!"

	var employeeToken string
	var customerToken string

	// --- Employee user setup ---
	t.Run("create employee", func(t *testing.T) {
		payload := map[string]string{
			"email":    employeeEmail,
			"password": password,
		}
		body, _ := json.Marshal(payload)
		resp, err := http.Post(baseURL+"/users/create-employee", "application/json", bytes.NewReader(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("login as employee", func(t *testing.T) {
		payload := map[string]string{
			"email":    employeeEmail,
			"password": password,
		}
		body, _ := json.Marshal(payload)
		resp, err := http.Post(baseURL+"/users/login", "application/json", bytes.NewReader(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]string
		_ = json.NewDecoder(resp.Body).Decode(&res)
		employeeToken = res["token"]
		assert.NotEmpty(t, employeeToken)
	})

	// --- Normal user setup ---
	t.Run("register customer", func(t *testing.T) {
		payload := map[string]string{
			"email":         customerEmail,
			"password":      password,
			"first_name":    "John",
			"last_name":     "Doe",
			"address_line1": "123 Main St",
			"city":          "Sortville",
			"postal_code":   "00000",
			"country":       "Sortland",
			"phone_number":  "1234567890",
		}
		body, _ := json.Marshal(payload)
		resp, err := http.Post(baseURL+"/users/register", "application/json", bytes.NewReader(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("login as customer", func(t *testing.T) {
		payload := map[string]string{
			"email":    customerEmail,
			"password": password,
		}
		body, _ := json.Marshal(payload)
		resp, err := http.Post(baseURL+"/users/login", "application/json", bytes.NewReader(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]string
		_ = json.NewDecoder(resp.Body).Decode(&res)
		customerToken = res["token"]
		assert.NotEmpty(t, customerToken)
	})

	// --- Employee access should work ---
	t.Run("employee can access special endpoint", func(t *testing.T) {
		req, _ := http.NewRequest("POST", baseURL+"/users/special", nil)
		req.Header.Set("Authorization", "Bearer "+employeeToken)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "Sort employee command executed successfully", res["message"])
		assert.Equal(t, "employee", res["role"])
	})

	// --- Normal user access should be denied ---
	t.Run("customer gets 403 on special endpoint", func(t *testing.T) {
		req, _ := http.NewRequest("POST", baseURL+"/users/special", nil)
		req.Header.Set("Authorization", "Bearer "+customerToken)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}
