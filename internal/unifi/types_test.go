package unifi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestSensitiveStringRedaction(t *testing.T) {
	s := SensitiveString("my-secret-api-key")

	cases := []struct {
		name string
		got  string
	}{
		{"%v", fmt.Sprintf("%v", s)}, //nolint:gocritic // testing fmt routing, not calling String() directly
		{"%s", fmt.Sprintf("%s", s)}, //nolint:gocritic // testing fmt routing, not calling String() directly
		{"%q", fmt.Sprintf("%q", s)},
		{"%#v", fmt.Sprintf("%#v", s)},
		{"%x", fmt.Sprintf("%x", s)},
		{"%X", fmt.Sprintf("%X", s)},
		{"String()", s.String()},
		{"GoString()", s.GoString()},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got == "my-secret-api-key" {
				t.Errorf("%s leaked the secret value", tc.name)
			}
		})
	}

	t.Run("MarshalJSON redacts", func(t *testing.T) {
		b, err := json.Marshal(s)
		if err != nil {
			t.Fatalf("MarshalJSON: %v", err)
		}
		if string(b) == `"my-secret-api-key"` {
			t.Error("MarshalJSON leaked the secret value")
		}
	})

	t.Run("string() exposes real value", func(t *testing.T) {
		if got := string(s); got != "my-secret-api-key" {
			t.Errorf("string(s) = %q, want real value", got)
		}
	})
}

func TestAPIError(t *testing.T) {
	t.Run("errors.As unwraps to APIError", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprint(w, "invalid credentials")
		})
		_, err := client.GetInfo(t.Context())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected *APIError via errors.As, got %T: %v", err, err)
		}
		if apiErr.StatusCode != http.StatusUnauthorized {
			t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusUnauthorized)
		}
		if apiErr.Body == "" {
			t.Error("expected non-empty Body")
		}
	})

	t.Run("large body is truncated", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			for range 10 {
				_, _ = fmt.Fprint(w, "x")
				for range 100 {
					_, _ = fmt.Fprint(w, "0123456789")
				}
			}
		})
		_, err := client.GetInfo(t.Context())
		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected *APIError, got %T", err)
		}
		if len(apiErr.Body) > maxErrBodyBytes+len("… (truncated)") {
			t.Errorf("Body len %d exceeds truncation limit", len(apiErr.Body))
		}
	})
}
