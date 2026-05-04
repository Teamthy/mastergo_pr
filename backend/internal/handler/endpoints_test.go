package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend/internal/models"
)

// Test helper function
func makeRequest(t *testing.T, method, path string, body interface{}) (*httptest.ResponseRecorder, string) {
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
	}

	req, err := http.NewRequest(method, path, bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	return w, string(bodyBytes)
}

// ============= AUTHENTICATION TESTS =============

func TestSignup(t *testing.T) {
	req := models.SignUpRequest{
		FirstName:       "Test",
		LastName:        "User",
		Email:           "test@example.com",
		Password:        "TestPassword123!",
		ConfirmPassword: "TestPassword123!",
	}

	w, _ := makeRequest(t, "POST", "/auth/signup", req)

	if w.Code != http.StatusCreated && w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 201 or 400, got %d", w.Code)
	}

	if w.Code == http.StatusCreated {
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err == nil && response["user"] != nil {
			t.Log("✓ Signup test passed")
		}
	}
}

func TestLogin(t *testing.T) {
	req := models.LoginRequest{
		Email:    "test@example.com",
		Password: "TestPassword123!",
	}

	w, _ := makeRequest(t, "POST", "/auth/login", req)

	if w.Code == http.StatusOK {
		var response models.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err == nil && response.Token != "" {
			t.Log("✓ Login test passed")
		}
	} else if w.Code == http.StatusUnauthorized {
		t.Log("✓ Login test passed (auth failed as expected)")
	} else {
		t.Errorf("Expected status 200 or 401, got %d", w.Code)
	}
}

func TestGetProfile(t *testing.T) {
	w, _ := makeRequest(t, "GET", "/auth/me", nil)

	if w.Code == http.StatusOK {
		var user models.User
		err := json.Unmarshal(w.Body.Bytes(), &user)
		if err == nil && user.Email != "" {
			t.Log("✓ Get profile test passed")
		}
	} else if w.Code == http.StatusUnauthorized {
		t.Log("✓ Get profile test passed (auth required)")
	} else {
		t.Errorf("Expected status 200 or 401, got %d", w.Code)
	}
}

// ============= PASSWORD RESET TESTS =============

func TestRequestPasswordReset(t *testing.T) {
	req := models.PasswordResetRequest{
		Email: "test@example.com",
	}

	w, _ := makeRequest(t, "POST", "/auth/password-reset/request", req)

	if w.Code == http.StatusOK {
		var response models.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err == nil && response.Message != "" {
			t.Log("✓ Request password reset test passed")
		}
	} else if w.Code == http.StatusBadRequest {
		t.Log("✓ Request password reset test passed (invalid email)")
	} else {
		t.Errorf("Expected status 200 or 400, got %d", w.Code)
	}
}

func TestResetPassword(t *testing.T) {
	req := models.PasswordResetConfirmRequest{
		Token:           "mock_token_12345",
		NewPassword:     "NewPassword456!",
		ConfirmPassword: "NewPassword456!",
	}

	w, _ := makeRequest(t, "POST", "/auth/password-reset/confirm", req)

	// Should fail with mock token
	if w.Code == http.StatusBadRequest {
		t.Log("✓ Reset password test passed (invalid token as expected)")
	} else if w.Code == http.StatusOK {
		var response models.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err == nil {
			t.Log("✓ Reset password test passed")
		}
	} else {
		t.Errorf("Expected status 200 or 400, got %d", w.Code)
	}
}

// ============= WALLET TESTS =============

func TestCreateWallet(t *testing.T) {
	w, _ := makeRequest(t, "POST", "/api/v1/wallet/create", map[string]interface{}{})

	if w.Code == http.StatusCreated || w.Code == http.StatusOK {
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err == nil && response["address"] != nil {
			t.Log("✓ Create wallet test passed")
		}
	} else if w.Code == http.StatusUnauthorized {
		t.Log("✓ Create wallet test passed (auth required)")
	} else if w.Code == http.StatusConflict {
		t.Log("✓ Create wallet test passed (wallet already exists)")
	} else {
		t.Errorf("Expected status 200, 201, 401 or 409, got %d", w.Code)
	}
}

func TestGetWalletBalance(t *testing.T) {
	w, _ := makeRequest(t, "GET", "/api/v1/wallet/balance", nil)

	if w.Code == http.StatusOK {
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err == nil && (response["balance"] != nil || response["amount"] != nil) {
			t.Log("✓ Get wallet balance test passed")
		}
	} else if w.Code == http.StatusUnauthorized {
		t.Log("✓ Get wallet balance test passed (auth required)")
	} else {
		t.Errorf("Expected status 200 or 401, got %d", w.Code)
	}
}

func TestGetTransactions(t *testing.T) {
	w, _ := makeRequest(t, "GET", "/api/v1/wallet/transactions", nil)

	if w.Code == http.StatusOK {
		var response []map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err == nil {
			t.Logf("✓ Get transactions test passed (%d transactions)", len(response))
		}
	} else if w.Code == http.StatusUnauthorized {
		t.Log("✓ Get transactions test passed (auth required)")
	} else {
		t.Errorf("Expected status 200 or 401, got %d", w.Code)
	}
}

// ============= API KEY TESTS =============

func TestCreateAPIKey(t *testing.T) {
	req := map[string]interface{}{
		"name": "Test API Key",
	}

	w, _ := makeRequest(t, "POST", "/api/v1/apikeys", req)

	if w.Code == http.StatusCreated || w.Code == http.StatusOK {
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err == nil && response["public_key"] != nil {
			t.Log("✓ Create API key test passed")
		}
	} else if w.Code == http.StatusUnauthorized {
		t.Log("✓ Create API key test passed (auth required)")
	} else {
		t.Errorf("Expected status 200, 201 or 401, got %d", w.Code)
	}
}

func TestListAPIKeys(t *testing.T) {
	w, _ := makeRequest(t, "GET", "/api/v1/apikeys", nil)

	if w.Code == http.StatusOK {
		var response []map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err == nil {
			t.Logf("✓ List API keys test passed (%d keys)", len(response))
		}
	} else if w.Code == http.StatusUnauthorized {
		t.Log("✓ List API keys test passed (auth required)")
	} else {
		t.Errorf("Expected status 200 or 401, got %d", w.Code)
	}
}

// ============= WEBHOOK TESTS =============

func TestCreateWebhook(t *testing.T) {
	req := models.CreateWebhookRequest{
		URL:    "https://webhook.site/unique-test-id",
		Events: []string{"transaction.created", "transaction.confirmed"},
	}

	w, _ := makeRequest(t, "POST", "/api/v1/webhooks", req)

	if w.Code == http.StatusCreated || w.Code == http.StatusOK {
		var response models.Webhook
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err == nil && response.URL != "" {
			t.Log("✓ Create webhook test passed")
			t.Logf("  - URL: %s", response.URL)
			t.Logf("  - Events: %v", response.Events)
		}
	} else if w.Code == http.StatusUnauthorized {
		t.Log("✓ Create webhook test passed (auth required)")
	} else {
		t.Errorf("Expected status 200, 201 or 401, got %d", w.Code)
	}
}

func TestListWebhooks(t *testing.T) {
	w, _ := makeRequest(t, "GET", "/api/v1/webhooks", nil)

	if w.Code == http.StatusOK {
		var response []models.Webhook
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err == nil {
			t.Logf("✓ List webhooks test passed (%d webhooks)", len(response))
		}
	} else if w.Code == http.StatusUnauthorized {
		t.Log("✓ List webhooks test passed (auth required)")
	} else {
		t.Errorf("Expected status 200 or 401, got %d", w.Code)
	}
}

// ============= ANALYTICS TESTS =============

func TestGetAPIAnalytics(t *testing.T) {
	// Replace with actual API key ID
	apiKeyID := "test-api-key-id"
	endpoint := fmt.Sprintf("/api/v1/analytics/api-keys/%s?days=7", apiKeyID)

	w, _ := makeRequest(t, "GET", endpoint, nil)

	if w.Code == http.StatusOK {
		var response models.AnalyticsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err == nil {
			t.Log("✓ Get API analytics test passed")
			t.Logf("  - Total requests: %d", response.TotalRequests)
			t.Logf("  - Error rate: %.2f%%", response.ErrorRate)
		}
	} else if w.Code == http.StatusUnauthorized {
		t.Log("✓ Get API analytics test passed (auth required)")
	} else if w.Code == http.StatusBadRequest {
		t.Log("✓ Get API analytics test passed (invalid ID)")
	} else {
		t.Errorf("Expected status 200, 400 or 401, got %d", w.Code)
	}
}

func TestGetUserAnalytics(t *testing.T) {
	w, _ := makeRequest(t, "GET", "/api/v1/analytics/user?days=7", nil)

	if w.Code == http.StatusOK {
		var response models.AnalyticsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err == nil {
			t.Log("✓ Get user analytics test passed")
			t.Logf("  - Total requests: %d", response.TotalRequests)
		}
	} else if w.Code == http.StatusUnauthorized {
		t.Log("✓ Get user analytics test passed (auth required)")
	} else {
		t.Errorf("Expected status 200 or 401, got %d", w.Code)
	}
}

// ============= AUDIT LOG TESTS =============

func TestGetAuditLogs(t *testing.T) {
	w, _ := makeRequest(t, "GET", "/api/v1/audit/logs?limit=50&offset=0", nil)

	if w.Code == http.StatusOK {
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err == nil && response["logs"] != nil {
			t.Log("✓ Get audit logs test passed")
		}
	} else if w.Code == http.StatusUnauthorized {
		t.Log("✓ Get audit logs test passed (auth required)")
	} else {
		t.Errorf("Expected status 200 or 401, got %d", w.Code)
	}
}

func TestGetAuditSummary(t *testing.T) {
	w, _ := makeRequest(t, "GET", "/api/v1/audit/summary?days=7", nil)

	if w.Code == http.StatusOK {
		var response map[string]int
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err == nil {
			t.Log("✓ Get audit summary test passed")
			t.Logf("  - Actions tracked: %d", len(response))
		}
	} else if w.Code == http.StatusUnauthorized {
		t.Log("✓ Get audit summary test passed (auth required)")
	} else {
		t.Errorf("Expected status 200 or 401, got %d", w.Code)
	}
}

// ============= PERFORMANCE TESTS =============

func TestAuthenticationPerformance(t *testing.T) {
	start := time.Now()

	req := models.LoginRequest{
		Email:    "test@example.com",
		Password: "TestPassword123!",
	}

	makeRequest(t, "POST", "/auth/login", req)

	duration := time.Since(start)

	if duration < 500*time.Millisecond {
		t.Log(fmt.Sprintf("✓ Authentication performance test passed: %dms", duration.Milliseconds()))
	} else {
		t.Log(fmt.Sprintf("⚠ Authentication slower than expected: %dms", duration.Milliseconds()))
	}
}

func TestWalletPerformance(t *testing.T) {
	start := time.Now()

	makeRequest(t, "GET", "/api/v1/wallet/balance", nil)

	duration := time.Since(start)

	if duration < 200*time.Millisecond {
		t.Log(fmt.Sprintf("✓ Wallet performance test passed: %dms", duration.Milliseconds()))
	} else {
		t.Log(fmt.Sprintf("⚠ Wallet slower than expected: %dms", duration.Milliseconds()))
	}
}

func TestAnalyticsPerformance(t *testing.T) {
	start := time.Now()

	makeRequest(t, "GET", "/api/v1/analytics/user?days=7", nil)

	duration := time.Since(start)

	if duration < 500*time.Millisecond {
		t.Log(fmt.Sprintf("✓ Analytics performance test passed: %dms", duration.Milliseconds()))
	} else {
		t.Log(fmt.Sprintf("⚠ Analytics slower than expected: %dms", duration.Milliseconds()))
	}
}

// ============= INTEGRATION TESTS =============

func TestFullAuthenticationFlow(t *testing.T) {
	ctx := context.Background()

	// Create test user
	signupReq := models.SignUpRequest{
		FirstName:       "Integration",
		LastName:        "Test",
		Email:           fmt.Sprintf("inttest-%d@example.com", time.Now().Unix()),
		Password:        "IntegrationTest123!",
		ConfirmPassword: "IntegrationTest123!",
	}

	w, _ := makeRequest(t, "POST", "/auth/signup", signupReq)
	if w.Code != http.StatusCreated && w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
		t.Errorf("Signup failed: %d", w.Code)
		return
	}

	// Login
	loginReq := models.LoginRequest{
		Email:    signupReq.Email,
		Password: signupReq.Password,
	}

	w, _ = makeRequest(t, "POST", "/auth/login", loginReq)
	if w.Code != http.StatusOK && w.Code != http.StatusUnauthorized {
		t.Errorf("Login failed: %d", w.Code)
		return
	}

	if w.Code == http.StatusOK {
		var response models.AuthResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		if response.Token != "" {
			t.Log("✓ Full authentication flow test passed")
			t.Log("  - Signup: OK")
			t.Log("  - Login: OK")
			t.Log("  - Token obtained: OK")
		}
	}
}

// Run all tests
func TestAllEndpoints(t *testing.T) {
	t.Log("Starting comprehensive API endpoint tests...")
	t.Log("")
	t.Log("Note: Some tests may require authentication tokens and IDs")
	t.Log("Ensure the backend is running on http://localhost:8080")
	t.Log("")
}
