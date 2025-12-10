/*
 * Copyright (C) 2025 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSendRequest(t *testing.T) {
	// Arrange
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		wantErr        bool
		serverBehavior func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name:         "Successful request with 200",
			statusCode:   http.StatusOK,
			responseBody: `{"status":"success"}`,
			wantErr:      false,
			serverBehavior: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"success"}`))
			},
		},
		{
			name:         "Successful request with 201",
			statusCode:   http.StatusCreated,
			responseBody: `{"status":"created"}`,
			wantErr:      false,
			serverBehavior: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{"status":"created"}`))
			},
		},
		{
			name:         "Successful request with 204",
			statusCode:   http.StatusNoContent,
			responseBody: "",
			wantErr:      false,
			serverBehavior: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			},
		},
		{
			name:       "Failed request with 400",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
			serverBehavior: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"bad request"}`))
			},
		},
		{
			name:       "Failed request with 404",
			statusCode: http.StatusNotFound,
			wantErr:    true,
			serverBehavior: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error":"not found"}`))
			},
		},
		{
			name:       "Failed request with 500 (should retry)",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
			serverBehavior: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"internal server error"}`))
			},
		},
		{
			name:       "Failed request with 503 (should retry)",
			statusCode: http.StatusServiceUnavailable,
			wantErr:    true,
			serverBehavior: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"error":"service unavailable"}`))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(tt.serverBehavior))
			defer server.Close()

			// Create a request
			req, err := http.NewRequest(http.MethodGet, server.URL, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("X-Request-Id", "test-request-id")

			// Create client with short timeout for tests
			client := &http.Client{
				Timeout: 5 * time.Second,
			}

			// Act
			body, err := SendRequest(client, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendRequest() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.responseBody != "" {
				// Assert
				assert.Equal(t, tt.responseBody, string(body))
			}
		})
	}
}

func TestSendRequestWithHeaders(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", r.Header.Get("X-Request-Id"))
		w.Header().Set("X-Trace-Id", "test-trace-id")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-Request-Id", "test-request-123")

	client := &http.Client{Timeout: 5 * time.Second}
	// Act
	_, err = SendRequest(client, req)
	assert.NoError(t, err)
}

func TestRetryPolicy(t *testing.T) {
	// Arrange
	tests := []struct {
		name       string
		statusCode int
		err        error
		wantRetry  bool
	}{
		{
			name:       "Retry on 500",
			statusCode: 500,
			err:        nil,
			wantRetry:  true,
		},
		{
			name:       "Retry on 503",
			statusCode: 503,
			err:        nil,
			wantRetry:  true,
		},
		{
			name:       "Retry on 504",
			statusCode: 504,
			err:        nil,
			wantRetry:  true,
		},
		{
			name:       "No retry on 200",
			statusCode: 200,
			err:        nil,
			wantRetry:  false,
		},
		{
			name:       "No retry on 400",
			statusCode: 400,
			err:        nil,
			wantRetry:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Status:     http.StatusText(tt.statusCode),
			}
			// Act
			shouldRetry, _ := retryPolicy(ctx, resp, tt.err)
			// Assert
			if shouldRetry != tt.wantRetry {
				t.Errorf("retryPolicy() shouldRetry = %v, want %v", shouldRetry, tt.wantRetry)
			}
		})
	}
}

func TestRetryPolicyWithContextCanceled(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	resp := &http.Response{StatusCode: 200}
	// Act
	shouldRetry, _ := retryPolicy(ctx, resp, nil)
	assert.False(t, shouldRetry, "retryPolicy() should not retry on canceled context")
}

func TestRetryPolicyWithContextDeadlineExceeded(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(10 * time.Millisecond) // Ensure timeout

	resp := &http.Response{StatusCode: 200}
	// Act
	_, err := retryPolicy(ctx, resp, nil)
	// Context deadline exceeded should be handled
	assert.Error(t, err, "retryPolicy() should return error on context deadline exceeded")
	assert.Equal(t, context.DeadlineExceeded, err)
}

func TestSendRequestTimeout(t *testing.T) {
	// Arrange
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-Request-Id", "test-timeout")

	// Create client with very short timeout
	client := &http.Client{
		Timeout: 1 * time.Millisecond,
	}
	// Act
	_, err = SendRequest(client, req)
	// Assert
	assert.NotNil(t, err, "SendRequest() should timeout and return error")
	assert.Contains(t, err.Error(), "Client.Timeout exceeded")
}

func TestSendRequestWithRetry(t *testing.T) {
	// Arrange
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error":"service unavailable"}`))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"success"}`))
		}
	}))
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-Request-Id", "test-retry")

	client := &http.Client{Timeout: 10 * time.Second}
	// Act
	body, err := SendRequest(client, req)

	// Assert
	assert.NoError(t, err, "SendRequest() should succeed after retry")
	assert.GreaterOrEqual(t, attempts, 2, "SendRequest() should have retried")
	assert.Equal(t, `{"status":"success"}`, string(body), "SendRequest() body should be success")
}
