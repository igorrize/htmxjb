package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	Db *sql.DB
}

func NewStore(dbName string) (Store, error) {
	dbPath := filepath.Join("data", dbName)
	Db, err := getConnection(dbPath)
	if err != nil {
		return Store{}, err
	}

	if err := createMigrations(dbName, Db); err != nil {
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
		// log.Fatalf("🔥 failed to connect to the database: %s", err.Error())
		return nil, fmt.Errorf("🔥 failed to connect to the database: %s", err)
	}

	log.Println("🚀 Connected Successfully to the Database")

	return db, nil
}

func createMigrations(dbName string, db *sql.DB) error {
	stmt := `CREATE TABLE IF NOT EXISTS jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_by INTEGER NOT NULL,
		title VARCHAR(64) NOT NULL,
		description VARCHAR(255) NULL,
		created_at DATETIME default CURRENT_TIMESTAMP
	);`

  _, err := db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}
