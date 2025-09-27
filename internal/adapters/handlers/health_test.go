package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MatheusHenrique129/bemax-api/internal/adapters/handlers"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	testCases := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "success",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			healthHandler := handlers.NewHealthHandler()

			url := "/ping"

			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()

			healthHandler.Ping(rec, req)
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
