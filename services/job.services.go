package services

import (
	"fmt"
	"github.com/igorrize/htmxjb/db"
	"time"
)

type Job struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description,omitempty"`
    CreatedAt   time.Time `json:"created_at,omitempty"`
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
	rows, err := js.JobStore.Db.Query(query)
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
