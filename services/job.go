package services

import (
	"fmt"
	"htmxjb/db"
	"htmxjb/models/domain"
	"time"
)

type Job struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	Type        string    `json:"type"`
	Location    string    `json:"location"`
	Salary      string    `json:"salary"`
	IsNew       bool      `json:"is_new"`
	Company     string    `json:"company"`
	Tags        []string  `json:"tags,omitempty"`
}

type JobServices struct {
	Job      Job
	JobStore db.Store
}

func NewJobServices(j Job, jobStore db.Store) *JobServices {
	return &JobServices{
		Job:      j,
		JobStore: jobStore,
	}
}

func (js *JobServices) GetAllJobs() ([]Job, error) {
	query := "SELECT id, title, description FROM jobs ORDER BY created_at DESC"
	rows, err := js.JobStore.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs: %w", err)
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
		if err := rows.Scan(&job.ID, &job.Title, &job.Description); err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}
		jobs = append(jobs, job)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating jobs: %w", err)
	}

	return jobs, nil
}

func (js *JobServices) Create(job *domain.Job) error {
	query := `
    INSERT INTO jobs (external_id, title, description, type, source)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id
  `
	var id int64
	err := js.JobStore.QueryRow(
		query,
		job.ExternalID,
		job.Title,
		job.Description,
		job.Type,
		job.Source,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	job.ID = id

	return nil
}
