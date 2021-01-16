package main

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/honeycombio/beeline-go/wrappers/hnysql"
	_ "github.com/lib/pq"
)

// Storage - Persistence layer
type Storage struct {
	database *hnysql.DB
}

// NewStorage - Creates a new Storage struct
func NewStorage(connStr string) (Storage, error) {
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return Storage{}, err
	}

	storage := Storage{database: hnysql.WrapDB(db)}

	return storage, nil
}

// PageExists - Checks if a page exists
func (s *Storage) PageExists(tenantID uuid.UUID, pageName string) (bool, error) {
	var exists bool
	cmd := "select exists(select 1 from tenant_page where tenant_id=$1 and name=$2)"

	err := s.database.
		QueryRow(cmd, tenantID, pageName).
		Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// InsertHit - Inserts a Hit record
func (s *Storage) InsertHit(tenantID uuid.UUID, pageName string, hitDate time.Time, footprint string) error {
	cmd := "insert into hit (tenant_id, page_name, event_time, footprint) values ($1, $2, $3, $4)"
	_, err := s.database.Exec(cmd, tenantID, pageName, hitDate, footprint)

	return err
}

// Disconnect - Closes the database connection
func (s *Storage) Disconnect() {
	if s.database == nil {
		return
	}

	s.database.Close()
}
