package store

import (
	"context"
	"sync"
	"errors"
	"database/sql"
	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	mu sync.Mutex	
	db *sql.DB
}

func NewSQLiteStore(filename string) (*SQLiteStore, error) {
	s := &SQLiteStore {}
	var err error
	s.db, err = sql.Open("sqlite", filename)

	if (err != nil) {
		return nil, err	
	}
	
	s.db.SetMaxOpenConns(1)

	_, err = s.db.Exec(`CREATE TABLE IF NOT EXISTS snippets (
		slug TEXT PRIMARY KEY,
		content TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`)

	if (err != nil) {
		return nil, err
	}

	return s, nil
}

func (s *SQLiteStore)  Save(ctx context.Context, slug string, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	result, err := s.db.ExecContext(ctx, "INSERT INTO snippets (slug, content) VALUES (?, ?)", slug, content)
	
	if (err != nil) {
		return err
	}

	rows, err := result.RowsAffected()

	if (err != nil) {
		return err
	}

	if (rows != 1) {
		return errors.New("Search came back empty")
	}

	return nil
}

func (s *SQLiteStore) Get(ctx context.Context, slug string) (string, error) {
	row := s.db.QueryRowContext(ctx, "SELECT slug, content FROM snippets WHERE slug=(?)")
	var slg, con string
	err := row.Scan(&slg, &con)
	if (err != nil) {
		return "", err
	}

	if (slg != slug) {
		return "", errors.New("Slug found doesn't match searched slug")
	}
	
	return con, nil
}

func (s *SQLiteStore) Delete(ctx context.Context, slug string) error {

	return nil 
}
