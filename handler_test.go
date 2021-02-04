package v8gohttp

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	http.Handle("/test", Handler(""))

	req := httptest.NewRequest("POST", "https://example.com/test", strings.NewReader(`{
		"name": "Dog 🐶"
	}`))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	http.DefaultServeMux.ServeHTTP(rec, req)

	res := rec.Result()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("res.StatusCode != 200, got %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	if string(body) != "Hello Dog 🐶!" {
		t.Fatalf("res.Body != \"Hello Dog 🐶!\", got %#v", string(body))
	}
}
