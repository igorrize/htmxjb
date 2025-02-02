package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
)

type Store struct {
	Db *sql.DB
}

func NewStore(dbName string) (Store, error) {
	if err := os.MkdirAll("data", 0755); err != nil {
		return Store{}, fmt.Errorf("failed to create data directory: %w", err)
	}

	dbPath := filepath.Join("data", dbName)

	Db, err := getConnection(dbPath)
	if err != nil {
		return Store{}, err
	}

	if err := createMigrations(dbPath, Db); err != nil {
		return Store{}, err
	}

	return Store{
		Db,
	}, nil
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
					created_by INTEGER NOT NULL,
					title VARCHAR(64) NOT NULL,
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
