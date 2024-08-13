package store

import (
	"database/sql"
	"fmt"
	"github.com/orainmers/golangStudy/internal/models"
	"log/slog"
	"net/url"
)

type Store struct {
	lg *slog.Logger
	db *sql.DB
}

func New(
	lg *slog.Logger,
	username string,
	password string,
	address string,
	database string,
) (*Store, error) {
	dsn := (&url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(username, password),
		Host:   address,
		Path:   database,
	}).String()

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("init db: %v", err)
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %v", err)
	}

	return &Store{
		lg: lg,
		db: sqlDB,
	}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) DummyMigration() error {
	query := `CREATE TABLE IF NOT EXISTS People
(
  id          uuid PRIMARY KEY   NOT NULL,
  name        VARCHAR            NOT NULL,
  description VARCHAR,
  created_at   timestamp          NOT NULL,
  updated_at   timestamp          NOT NULL,
  is_deleted  bool DEFAULT FALSE NOT NULL
);`

	if _, err := s.db.Exec(query); err != nil {
		return fmt.Errorf("create table: %v", err)
	}

	return nil
}

func (s *Store) AddPerson(person *models.Person) error {
	query := `INSERT INTO People (id, name, description, created_at, updated_at, is_deleted)
VALUES ($1, $2, $3, $4, $5, $6);`

	if _, err := s.db.Exec(
		query,
		person.ID,
		person.Name,
		person.Description,
		person.CreatedAt,
		person.UpdatedAt,
		person.IsDeleted,
	); err != nil {
		return fmt.Errorf("s.db.Exec(...): %v", err)
	}

	return nil
}
