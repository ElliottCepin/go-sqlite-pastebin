package server

import (
	"context"
	"net/http"
	"log/slog"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
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

	if (r.Header.Get("Content-Type") != "plain/text") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	slug := generateSlug()
	encoder := json.NewEncoder(w)
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	err = encoder.Encode(
		Data{
			Slug: slug,
		})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	
	s.store.Save(context.Background(), slug, string(body))
	w.WriteHeader(http.StatusOK)
}

func generateSlug() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	token := hex.EncodeToString(b)
	return token[:8]	
}

func (s *Server) handleRetrieve() {

}

func (s *Server) handleDelete() {

}

func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /snippet", s.handleCreate)
	return mux
}
