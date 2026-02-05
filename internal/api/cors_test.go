package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithCORS_AddsHeaders(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	handler := WithCORS(next)
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("missing Access-Control-Allow-Origin")
	}
	if w.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Fatalf("missing Access-Control-Allow-Methods")
	}
	if w.Header().Get("Access-Control-Allow-Headers") == "" {
		t.Fatalf("missing Access-Control-Allow-Headers")
	}
	if !called {
		t.Fatalf("expected next handler to be called")
	}
}

func TestWithCORS_NonOptions_CallsNext(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	handler := WithCORS(next)
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if !called {
		t.Fatalf("expected next handler to be called")
	}
}
