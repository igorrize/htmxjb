package indeed_client

import (
	"encoding/json"
	"fmt"
	"htmxjb/models/responses"
	"htmxjb/clients/csv_client"
	"io"
	"net/http"
)

// Job представляет универсальную модель вакансии для сохранения в БД
type Job struct {
	ExternalID  string
	Title       string
	Description string
	Company     string
	Location    string
	URL         string
	Source      int // 1 для Indeed
	Type        string
}

// Интерфейс клиента для Indeed
type IndeedClientInterface interface {
	GetJobs() ([]Job, error)
	SaveJobs(db Database) error
}

// Реализация клиента
type IndeedClient struct {
	apiKey string
	host   string
}

// NewIndeedClient создает новый экземпляр клиента Indeed
func NewIndeedClient(apiKey string) *IndeedClient {
	return &IndeedClient{
		apiKey: apiKey,
		host:   "indeed-scraper-api.p.rapidapi.com",
	}
}

// GetJobs выполняет запрос к API Indeed и возвращает вакансии
func (c *IndeedClient) GetJobs() ([]Job, error) {
	url := "https://indeed-scraper-api.p.rapidapi.com/jobs"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Добавление параметров запроса
	q := req.URL.Query()
	q.Add("query", "htmx")

	req.URL.RawQuery = q.Encode()

	// Установка заголовков
	req.Header.Add("x-rapidapi-key", c.apiKey)
	req.Header.Add("x-rapidapi-host", c.host)

	// Выполнение запроса
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	// Проверка статус-кода
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("received non-200 response: %d, body: %s", res.StatusCode, string(body))
	}

	// Чтение и разбор ответа
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var response responses.IndeedResponse
	if err := json.Unmarshal(body, &response); err != nil {
		// Для отладки выводим часть тела ответа
		previewLen := 200
		if len(body) < previewLen {
			previewLen = len(body)
		}
		return nil, fmt.Errorf("error parsing JSON: %w, response preview: %s", err, string(body[:previewLen]))
	}

	// Преобразование данных из API в нашу структуру Job
	var jobs []Job
	for _, indeedJob := range response.ReturnValue.Data {
		location := ""
		if indeedJob.Location.FormattedAddressShort != "" {
			location = indeedJob.Location.FormattedAddressShort
		} else if indeedJob.Location.FormattedAddressLong != "" {
			location = indeedJob.Location.FormattedAddressLong
		} else if indeedJob.Location.City != "" {
			location = indeedJob.Location.City
		}

		jobs = append(jobs, Job{
			ExternalID:  indeedJob.JobKey,
			Title:       indeedJob.Title,
			Description: indeedJob.DescriptionText,
			Company:     indeedJob.CompanyName,
			Location:    location,
			URL:         indeedJob.JobUrl,
			Source:      1,
			Type:        indeedJob.JobType,
		})
	}

	return jobs, nil
}

// SaveJobs сохраняет вакансии в базу данных
func (c *IndeedClient) SaveJobs(db Database) error {
	jobs, err := c.GetJobs()
	if err != nil {
		return fmt.Errorf("error fetching jobs: %w", err)
	}

	for _, job := range jobs {
		// Проверяем существование вакансии перед сохранением
		existingJob, err := db.GetJobByExternalID(job.ExternalID)
		if err == nil && existingJob.ExternalID != "" {
			// Вакансия уже существует, пропускаем
			continue
		}

		// Сохраняем новую вакансию
		if err := db.SaveJob(job); err != nil {
			return fmt.Errorf("error saving job %s: %w", job.ExternalID, err)
		}
	}

	return nil
}

// Интерфейс для базы данных
type Database interface {
	SaveJob(job Job) error
	GetJobByExternalID(externalID string) (Job, error)
}
