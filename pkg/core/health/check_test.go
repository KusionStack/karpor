package health

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockCheck struct {
	name string
	pass bool
}

func (m *mockCheck) Name() string {
	return m.name
}

func (m *mockCheck) Pass(ctx context.Context) bool {
	return m.pass
}

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name               string
		query              string
		checks             []Check
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:               "All checks pass",
			query:              "",
			checks:             []Check{&mockCheck{name: "check1", pass: true}, &mockCheck{name: "check2", pass: true}},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `OK`,
		},
		{
			name:               "One check fails",
			query:              "",
			checks:             []Check{&mockCheck{name: "check1", pass: true}, &mockCheck{name: "check2", pass: false}},
			expectedStatusCode: http.StatusServiceUnavailable,
			expectedResponse:   `Fail`,
		},
		{
			name:               "Check excluded",
			query:              "excludes=check2",
			checks:             []Check{&mockCheck{name: "check1", pass: true}, &mockCheck{name: "check2", pass: false}},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `OK`,
		},
		{
			name:               "Verbose output",
			query:              "verbose=true",
			checks:             []Check{&mockCheck{name: "check1", pass: true}, &mockCheck{name: "check2", pass: true}},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `[+] check1 ok` + "\n" + `[+] check2 ok` + "\n" + `health check passed`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := HandlerConfig{
				Verbose:  false,
				Excludes: []string{},
				Checks:   tt.checks,
				FailureNotification: FailureNotification{
					Threshold: 3,
					Chan:      make(chan error),
				},
			}

			req, err := http.NewRequest("GET", "/health?"+tt.query, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := NewHandler(conf)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatusCode)
			}

			if response := strings.TrimSpace(rr.Body.String()); response != tt.expectedResponse {
				t.Errorf("handler returned unexpected body: got %v want %v",
					response, tt.expectedResponse)
			}
		})
	}
}
