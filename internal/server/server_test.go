package server

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/ElliottCepin/go-sqlite-pastebin/internal/store"
	"log/slog"
	"strings"
	"encoding/json"
	regex "regexp"
	"io"
)

func TestRoundTrip(t *testing.T) {
	s := &Server{
		store: store.NewMemoryStore(),
		logger: slog.Default(),
	}

	srv := httptest.NewServer(s.Routes())
	t.Cleanup(srv.Close)

	resp, err := http.Post(srv.URL + "/snippet", "text/plain", strings.NewReader("Content"))
	
	if err != nil {
		t.Errorf("Issue with POST /snippet: %v | status code: %v", err, resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	var data Data

	
	err = dec.Decode(&data)

	if err != nil {
		t.Errorf("Could not decode response: %v | status code: %v", err, resp.StatusCode)
	}

	slug := data.Slug

	re := "^[A-Za-z0-9]{8}$"
	rex, err := regex.Compile(re)

	if err != nil {
		t.Errorf("Issues with regex (%v): %v", re, err)
	}

	if (!rex.MatchString(slug)) {
		t.Errorf("Slug '%v' does not match re '%v'", slug, re)
	}

	resp, err = http.Get(srv.URL + "/snippet/" + slug)
	
	if (err != nil) {
		t.Errorf("GET %v/snippet/%v failed: %v", srv.URL, slug, err)
	}

	if (resp.StatusCode != http.StatusOK) {
		t.Errorf("GET %v/snippet/%v returned status: %v", srv.URL, slug, resp.StatusCode)
	}

	if (resp.Header.Get("Content-Type") != "text/plain") {
		t.Errorf("Expected Content-Type text/plain, got '%v'", resp.Header.Get("Content-Type"))
	}

	body, err := io.ReadAll(resp.Body) 

	if (err != nil) {
		t.Errorf("Could not read body: %v", err)
	}

	if (string(body) != "Content") {
		t.Errorf("Expected 'Content' got '%v'", string(body))
	}


}

func TestDelete(t *testing.T) {
	s := &Server{
		store: store.NewMemoryStore(),
		logger: slog.Default(),
	}

	srv := httptest.NewServer(s.Routes())

	resp, err := http.Post(srv.URL + "/snippet", "text/plain", strings.NewReader("Content"))
	
	if err != nil {
		t.Errorf("Issue with POST /snippet: %v | status code: %v", err, resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	var data Data

	
	err = dec.Decode(&data)

	if err != nil {
		t.Errorf("Could not decode response: %v | status code: %v", err, resp.StatusCode)
	}

	slug := data.Slug
	
	req, err := http.NewRequest(http.MethodDelete, srv.URL + "/snippet/" + slug, nil)

	if (err != nil) {
		t.Errorf("Error constructing request: %v", err) 
	}

	res, err := http.DefaultClient.Do(req)

	if (err != nil || res.StatusCode != http.StatusNoContent) {
		t.Errorf("Error issuing request: %v | status code: %v", err, res.StatusCode)
	}

	res, err = http.Get(srv.URL + "/snippet/" + slug)

	if (res.StatusCode != http.StatusNotFound) {
		t.Errorf("Expected status code %v, got: %v", http.StatusNotFound, res.StatusCode)
	}

	

}
