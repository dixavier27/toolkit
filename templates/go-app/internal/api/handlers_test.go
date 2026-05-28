package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHello_OK(t *testing.T) {
	mux := NewMux("test")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/hello/eco", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("esperava 200, recebi %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Olá, eco!") {
		t.Errorf("body inesperado: %s", rec.Body.String())
	}
}

func TestHealthz(t *testing.T) {
	mux := NewMux("v1.2.3")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("esperava 200, recebi %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "v1.2.3") {
		t.Errorf("esperava versão no body, recebi: %s", rec.Body.String())
	}
}
