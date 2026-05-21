package server

import (
	"context"
	"net/http"
	"slog"
	"crypto/rand"
	"hex"
)

type Store interface {
	Save(ctx context.Context, slug string, content string) error
	Get(ctx context.Context, slug string) (string, error)
	Delete(ctx context.Context, slug string) error
}

type Server interface {
	store Store
	logger *slog.Logger
}

func (s *Server) handleCreate(w http.ResponseWriter, r *http.Request){
		
}

func generateSlug() string {
	b := make([]byte, 32)
	_ := rand.Read(b)
	token := hex.EncodeToString(b)
	return token[:8]	
}

func (s *Server) handleRetrieve() {

}

func (s *Server) handleDelete() {

}

func (s *Server) Routes() *http.ServeMux {

}
