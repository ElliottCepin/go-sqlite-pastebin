package server

import (
	"context"
	"net/http"
	"log/slog"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"fmt"
	"mime"
)

type Store interface {
	Save(ctx context.Context, slug string, content string) error
	Get(ctx context.Context, slug string) (string, error)
	Delete(ctx context.Context, slug string) error
}

type Server struct {
	store Store
	logger *slog.Logger
}

type Data struct {
	Slug string
}

func (s *Server) handleCreate(w http.ResponseWriter, r *http.Request){
	if (r.Method != "POST") {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if (err != nil || mt != "text/plain") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	slug := generateSlug()
	encoder := json.NewEncoder(w)
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = encoder.Encode(
		Data{
			Slug: slug,
		})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}


	
	s.store.Save(context.Background(), slug, string(body))
}

func generateSlug() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	token := hex.EncodeToString(b)
	return token[:8]	
}

func (s *Server) handleRetrieve(w http.ResponseWriter, r *http.Request) {
	if (r.Method != "GET") {
		w.WriteHeader(http.StatusMethodNotAllowed)	
		return
	}

	slug := r.PathValue("slug")
	content, err := s.store.Get(context.Background(), slug)

	if (err != nil) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "%v", content)
	
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	if (r.Method != "DELETE") {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	slug := r.PathValue("slug")
	err := s.store.Delete(context.Background(), slug)

	if (err != nil) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /snippet", s.handleCreate)
	mux.HandleFunc("GET /snippet/{slug}", s.handleRetrieve)
	mux.HandleFunc("DELETE /snippet/{slug}", s.handleDelete)
	return mux
}
