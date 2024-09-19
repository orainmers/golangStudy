package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/orainmers/golangStudy/internal/models"
	"log/slog"
	"net/url"
	"strings"
)

const moduleName = "storage"

type Storage struct {
	lg *slog.Logger
	db *sql.DB
}

func New(
	lg *slog.Logger,
	username string,
	password string,
	address string,
	database string,
) (*Storage, error) {
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

	return &Storage{
		lg: lg.With("module", moduleName),
		db: sqlDB,
	}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) DummyMigration() error {
	query := `CREATE TABLE IF NOT EXISTS contacts
(
    id          uuid PRIMARY KEY   NOT NULL DEFAULT gen_random_uuid(),
    name        VARCHAR            NOT NULL,
    description VARCHAR,
    created_at   timestamp          NOT NULL default now(),
    updated_at   timestamp          NOT NULL default now(),
    deleted  bool DEFAULT FALSE NOT NULL
);`

	if _, err := s.db.Exec(query); err != nil {
		return fmt.Errorf("create table: %v", err)
	}

	s.lg.Info("Migration is succeed...")

	return nil
}

func (s *Storage) AddContact(ctx context.Context, contact models.ContactRequest) (uuid.UUID, error) {
	query := `INSERT INTO contacts (name, description)
VALUES ($1, $2);`

	c := s.db.QueryRowContext(
		ctx,
		query,
		contact.Name,
		contact.Description,
	)
	if c == nil {
		return uuid.Nil, fmt.Errorf("s.db.QueryRowContext(...): %v", ErrInternal)
	}

	result := struct {
		id uuid.UUID `db:"id"`
	}{}

	if err := c.Scan(&result.id); err != nil {
		return uuid.Nil, fmt.Errorf("c.Scan(...): %v", err)
	}

	return result.id, nil
}

func (s *Storage) GetContact(ctx context.Context, id uuid.UUID) (models.Contact, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM contacts
WHERE id = $1 AND deleted = false;`

	c := s.db.QueryRowContext(ctx, query, id.String())
	if c == nil {
		return models.Contact{}, ErrNotFound
	}

	var contact models.Contact
	if err := c.Scan(&contact); err != nil {
		return models.Contact{}, fmt.Errorf("rows.Scan(...): %v", err)
	}

	if contact.ID == uuid.Nil {
		return models.Contact{}, ErrNotFound
	}

	return contact, nil
}

func (s *Storage) GetContacts(ctx context.Context) ([]models.Contact, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM contacts
	WHERE NOT deleted;`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("s.db.QueryContext: %v", err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			s.lg.Info("rows.Close()", "err", err.Error())
		}
	}()

	var contacts []models.Contact

	for rows.Next() {
		var contact models.Contact

		if err = rows.Scan(
			&contact.ID,
			&contact.Name,
			&contact.Description,
			&contact.CreatedAt,
			&contact.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan(...): %v", err)
		}

		contacts = append(contacts, contact)
	}

	if len(contacts) == 0 {
		return nil, ErrNotFound
	}

	return contacts, nil
}

func (s *Storage) UpdateContact(ctx context.Context, id uuid.UUID, contact models.ContactRequest) (models.Contact, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return models.Contact{}, fmt.Errorf("open transaction faild: %w", err)
	}
	defer func() {
		if err = tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			s.lg.Warn("rollback transaction", "err", err.Error())
		}
	}()

	var args []any
	var query strings.Builder

	query.WriteString(`UPDATE contacts SET` + ` `)

	if contact.Name != nil {
		args = append(args, *contact.Name)
		query.WriteString(`name = $` + fmt.Sprint(len(args)) + `, `)
	}

	if contact.Description != nil {
		args = append(args, *contact.Description)
		query.WriteString(`description = $` + fmt.Sprint(len(args)) + `, `)
	}

	args = append(args, id)
	query.WriteString(fmt.Sprintf(` updated_at = NOW() WHERE id = $%d
RETURNING id, name, description, created_at, updated_at;`, len(args)))

	c := tx.QueryRowContext(ctx, query.String(), args...)
	if c == nil {
		return models.Contact{}, ErrNotFound
	}

	var updatedContact models.Contact
	if err = c.Scan(&updatedContact); err != nil {
		return models.Contact{}, fmt.Errorf("c.Scan(...): %v", err)
	}

	return updatedContact, nil
}

func (s *Storage) DeleteContact(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE contacts
SET deleted = true,
    updated_at = now()
WHERE id = $1;`

	if _, err := s.db.ExecContext(ctx, query, id.String()); err != nil {
		return fmt.Errorf("s.db.ExecContext: %v", err)
	}

	return nil
}
