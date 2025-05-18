package services

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type MockStore struct {
	Db *sql.DB
}

func (m *MockStore) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return m.Db.Query(query, args...)
}

func (m *MockStore) QueryRow(query string, args ...interface{}) *sql.Row {
	return m.Db.QueryRow(query, args...)
}

func (m *MockStore) Close() error {
	return m.Db.Close()
}
func TestGetAllJobs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockStore := &MockStore{Db: db}

	jobServices := NewJobServices(Job{}, mockStore)

	t.Run("Successfully get all jobs", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "description"}).
			AddRow(1, "Software Engineer", "Develop software").
			AddRow(2, "Data Scientist", "Analyze data")

		mock.ExpectQuery("SELECT id, title, description FROM jobs ORDER BY created_at DESC").WillReturnRows(rows)

		jobs, err := jobServices.GetAllJobs()

		assert.NoError(t, err)
		assert.Equal(t, 2, len(jobs))
		assert.Equal(t, "Software Engineer", jobs[0].Title)
		assert.Equal(t, "Data Scientist", jobs[1].Title)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Handle database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, title, description FROM jobs ORDER BY created_at DESC").WillReturnError(fmt.Errorf("mock database error"))

		_, err := jobServices.GetAllJobs()

		assert.Error(t, err)
		assert.Equal(t, "failed to get jobs: mock database error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
