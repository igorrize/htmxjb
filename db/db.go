package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
)

type SQLiteStore struct {
	Db *sql.DB
}

type Store interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Close() error
}

func NewStore(dbName string) (SQLiteStore, error) {
	if err := os.MkdirAll("data", 0755); err != nil {
		return SQLiteStore{}, fmt.Errorf("failed to create data directory: %w", err)
	}

	dbPath := filepath.Join("data", dbName)

	Db, err := getConnection(dbPath)
	if err != nil {
		return SQLiteStore{}, err
	}

	if err := createMigrations(dbPath, Db); err != nil {
		return SQLiteStore{}, err
	}

	return SQLiteStore{
		Db,
	}, nil
}

func (s *SQLiteStore) Close() error {
	if s.Db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	return s.Db.Close()
}

func getConnection(dbName string) (*sql.DB, error) {
	var (
		err error
		db  *sql.DB
	)

	if db != nil {
		return db, nil
	}

	db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, fmt.Errorf("🔥 failed to connect to the database: %s", err)
	}

	log.Println("🚀 Connected Successfully to the Database")

	return db, nil
}

func createMigrations(dbName string, db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			executed_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	migrations := []struct {
		name string
		stmt string
	}{
		{
			name: "create_jobs_table",
			stmt: `
				CREATE TABLE IF NOT EXISTS jobs (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					type INTEGER NOT NULL,
					title VARCHAR(64) NOT NULL,
          source INTEGER NOT NULL,
          external_id INTEGER NOT NULL,
					description VARCHAR(255) NULL,
					created_at DATETIME default CURRENT_TIMESTAMP);`,
		},
		{
			name: "add_updated_at_column_to_jobs",
			stmt: `
				ALTER TABLE jobs ADD COLUMN updated_at DATETIME default CURRENT_TIMESTAMP;`,
		},
	}

	for _, migration := range migrations {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = ?", migration.name).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if count == 0 {
			_, err := db.Exec(migration.stmt)
			if err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", migration.name, err)
			}

			_, err = db.Exec("INSERT INTO migrations (name) VALUES (?)", migration.name)
			if err != nil {
				return fmt.Errorf("failed to record migration %s: %w", migration.name, err)
			}

			log.Printf("✅ Migration %s executed successfully", migration.name)
		}
	}

	return nil
}
