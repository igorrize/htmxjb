package indeed_client

import (
	"fmt"
	"io"
	"net/http"
)

type IndeedClient struct {
	apiKey string
	host   string
}

func NewIndeedClient(apiKey string) *IndeedClient {
	return &IndeedClient{
		apiKey: apiKey,
		host:   "indeed-scraper-api.p.rapidapi.com",
	}
}

func (c *IndeedClient) GetJobs() ([]domain, Job, error) {
	url := "https://indeed-scraper-api.p.rapidapi.com/api/job/gtk6rn0y35rcm5r0hz5qvf0i"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("x-rapidapi-key", c.apiKey)
	req.Header.Add("x-rapidapi-host", c.host)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	var indeedJobs []IndeedJob
	if err := json.Unmarshal(body, &IndeedJob); err != nil {
		return nil, err
	}

	var jobs []Job
	for _, ij := range IndeedJobs {
		jobs.append(jobs, Job{
			ExternalID:  ij.ID,
			Title:       ij.Title,
			Description: ij.Description,
			Source:      1,
			Type:        ij.Type,
		})
	}
	return jobs, nil
}
