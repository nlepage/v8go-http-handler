package v8gohttp_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	v8gohttp "github.com/nlepage/v8go-http-handler"
)

func TestHandle(t *testing.T) {
	v8gohttp.Handle("/test", `
		async function handler(e) {
			const { name } = await e.request.json()
			e.respondWith(new Response('Hello ' + name + '!'))
		}
	`)

	srv := httptest.NewServer(http.DefaultServeMux)
	defer srv.Close()

	res, err := srv.Client().Post(srv.URL+"/test", "application/json", strings.NewReader(`{
		"name": "Dog 🐶"
	}`))
	if err != nil {
		panic(err)
	}

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
