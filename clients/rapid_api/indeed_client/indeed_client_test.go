package indeed_client

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestNewIndeedClient(t *testing.T) {
	apiKey := "test-api-key"
	client := NewIndeedClient(apiKey)
	
	if client.apiKey != apiKey {
		t.Errorf("Expected apiKey to be %s, got %s", apiKey, client.apiKey)
	}
	
	expectedHost := "indeed-scraper-api.p.rapidapi.com"
	if client.host != expectedHost {
		t.Errorf("Expected host to be %s, got %s", expectedHost, client.host)
	}
}

// MockRoundTripper реализует интерфейс http.RoundTripper для тестирования
type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func TestGetJobs_Success(t *testing.T) {
	// Подготовка тестовых данных
	mockJobs := []IndeedJob{
		{
			ID:          "job1",
			Title:       "Software Engineer",
			Description: "Job description 1",
			Type:        "Full-time",
		},
		{
			ID:          "job2",
			Title:       "Data Scientist", 
			Description: "Job description 2",
			Type:        "Contract",
		},
	}
	
	// Преобразование в JSON
	mockJSON, err := json.Marshal(mockJobs)
	if err != nil {
		t.Fatalf("Failed to marshal mock jobs: %v", err)
	}
	
	// Создание мок-клиента HTTP
	mockClient := &http.Client{
		Transport: &MockRoundTripper{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Проверка заголовков
				apiKey := req.Header.Get("x-rapidapi-key")
				if apiKey != "test-api-key" {
					t.Errorf("Expected API key header to be test-api-key, got %s", apiKey)
				}
				
				host := req.Header.Get("x-rapidapi-host")
				if host != "indeed-scraper-api.p.rapidapi.com" {
					t.Errorf("Expected host header to be indeed-scraper-api.p.rapidapi.com, got %s", host)
				}
				
				// Формирование ответа
				response := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(string(mockJSON))),
					Header:     make(http.Header),
				}
				return response, nil
			},
		},
	}
	
	// Временная замена стандартного HTTP-клиента
	defaultClient := http.DefaultClient
	http.DefaultClient = mockClient
	defer func() { http.DefaultClient = defaultClient }()
	
	// Создание клиента Indeed и выполнение запроса
	client := NewIndeedClient("test-api-key")
	jobs, err := client.GetJobs()
	
	// Проверки
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if len(jobs) != 2 {
		t.Errorf("Expected 2 jobs, got %d", len(jobs))
	}
	
	// Проверка первой работы
	if jobs[0].ExternalID != "job1" {
		t.Errorf("Expected job ID job1, got %s", jobs[0].ExternalID)
	}
	if jobs[0].Title != "Software Engineer" {
		t.Errorf("Expected job title Software Engineer, got %s", jobs[0].Title)
	}
	
	// Проверка второй работы
	if jobs[1].ExternalID != "job2" {
		t.Errorf("Expected job ID job2, got %s", jobs[1].ExternalID)
	}
	if jobs[1].Type != "Contract" {
		t.Errorf("Expected job type Contract, got %s", jobs[1].Type)
	}
}

func TestGetJobs_HttpError(t *testing.T) {
	// Создание мок-клиента, возвращающего ошибку
	mockClient := &http.Client{
		Transport: &MockRoundTripper{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Возврат ошибки 500
				response := &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader("Internal Server Error")),
					Header:     make(http.Header),
				}
				return response, nil
			},
		},
	}
	
	// Временная замена стандартного HTTP-клиента
	defaultClient := http.DefaultClient
	http.DefaultClient = mockClient
	defer func() { http.DefaultClient = defaultClient }()
	
	// Выполнение теста
	client := NewIndeedClient("test-api-key")
	_, err := client.GetJobs()
	
	// Проверка наличия ошибки
	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

func TestGetJobs_InvalidJson(t *testing.T) {
	// Создание мок-клиента, возвращающего неверный JSON
	mockClient := &http.Client{
		Transport: &MockRoundTripper{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Возврат неверного JSON
				response := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("invalid json")),
					Header:     make(http.Header),
				}
				return response, nil
			},
		},
	}
	
	// Временная замена стандартного HTTP-клиента
	defaultClient := http.DefaultClient
	http.DefaultClient = mockClient
	defer func() { http.DefaultClient = defaultClient }()
	
	// Выполнение теста
	client := NewIndeedClient("test-api-key")
	_, err := client.GetJobs()
	
	// Проверка наличия ошибки при парсинге JSON
	if err == nil {
		t.Error("Expected a JSON parsing error, got nil")
	}
}
