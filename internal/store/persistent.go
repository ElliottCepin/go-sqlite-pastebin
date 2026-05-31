package store

import (
	"context"
	"errors"
	"database/sql"
	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
	del *sql.Stmt
	get *sql.Stmt
	set *sql.Stmt
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

	s.set, err = s.db.Prepare("INSERT INTO snippets (slug, content) VALUES (?, ?);")
	s.get, err = s.db.Prepare("SELECT slug, content FROM snippets WHERE slug=(?)")
	s.del, err = s.db.Prepare("DELETE FROM snippets WHERE slug=(?)")

	return s, nil
}

func (s *SQLiteStore)  Save(ctx context.Context, slug string, content string) error {

	result, err := s.set.ExecContext(ctx, slug, content)
	
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
	row := s.get.QueryRowContext(ctx, slug)
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
	_, err := s.del.ExecContext(ctx, slug)
	if (err != nil) {
		return err
	}

	return nil 
}

func (s *SQLiteStore) Close() {
	s.del.Close()
	s.get.Close()
	s.set.Close()
	s.db.Close()
}
