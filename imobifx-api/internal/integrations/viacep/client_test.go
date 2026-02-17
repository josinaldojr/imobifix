package viacep_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/josinaldojr/imobifix-api/internal/integrations/viacep"
)

func TestLookup_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"cep":"58000-000","logradouro":"Rua X","bairro":"Bairro Y","localidade":"João Pessoa","uf":"PB"}`))
	}))
	defer srv.Close()

	c := viacep.NewClient(srv.URL, 500*time.Millisecond)

	addr, err := c.Lookup(context.Background(), "58000000")
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if addr.City != "João Pessoa" || addr.State != "PB" {
		t.Fatalf("unexpected addr: %+v", addr)
	}
}

func TestLookup_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"erro": true}`))
	}))
	defer srv.Close()

	c := viacep.NewClient(srv.URL, 500*time.Millisecond)

	_, err := c.Lookup(context.Background(), "58000000")
	if err != viacep.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestLookup_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`not-json`))
	}))
	defer srv.Close()

	c := viacep.NewClient(srv.URL, 500*time.Millisecond)

	_, err := c.Lookup(context.Background(), "58000000")
	if err != viacep.ErrInvalidAnswer {
		t.Fatalf("expected ErrInvalidAnswer, got %v", err)
	}
}

func TestLookup_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(120 * time.Millisecond)
		w.WriteHeader(200)
		w.Write([]byte(`{"cep":"58000-000"}`))
	}))
	defer srv.Close()

	c := viacep.NewClient(srv.URL, 50*time.Millisecond)

	_, err := c.Lookup(context.Background(), "58000000")
	if err != viacep.ErrUnavailable {
		t.Fatalf("expected ErrUnavailable, got %v", err)
	}
}