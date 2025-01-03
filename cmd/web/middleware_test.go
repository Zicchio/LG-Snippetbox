package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// next is mock next-in-chain middleware
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	secureHeaders(next).ServeHTTP(rr, r)
	rs := rr.Result()
	if frameOpt := rs.Header.Get("X-Frame-Options"); frameOpt != "deny" {
		t.Errorf("want %q; got %q", "deny", frameOpt)
	}

	if xssProtection := rs.Header.Get("X-XSS-Protection"); xssProtection != "1; mode=block" {
		t.Errorf("want %q; got %q", "1; mode=block", xssProtection)
	}

	if obt := rs.StatusCode; obt != http.StatusOK {
		t.Errorf("want %q; obt %q", http.StatusOK, obt)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	if obt := string(body); obt != "OK" {
		t.Errorf("want %q; got %q", "OK", obt)
	}
}
