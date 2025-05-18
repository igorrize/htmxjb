package csv_client

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
)

type CSVparser interface {
	ParseCSV(filePath string, hasHeader bool, numWorkers int) ([]Job, error)
}

func ParseCSV(filePath string, hasHeader bool, numWorkers int) ([]Job, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("ERROR while open file: %v", err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)

	if hasHeader {
		if _, err := csvReader.Read(); err != nil {
			return nil, fmt.Errorf("ERROR while read headers CSV: %v", err)
		}
	}

	rowsChan := make(chan []string)
	jobsChan := make(chan Job)
	errChan := make(chan error)
	doneChan := make(chan struct{})

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processRows(rowsChan, jobsChan, errChan, &wg)
	}

	go func() {
		for {
			record, err := csvReader.Read()
			if err != nil {
				if err.Error() == "EOF" {
					close(rowsChan)
					return
				}
				errChan <- fmt.Errorf("ERROR while reading CSV: %v", err)
				return
			}
			rowsChan <- record
		}
	}()

	var jobs []Job

	go func() {
		for job := range jobsChan {
			jobs = append(jobs, job)
		}
		close(doneChan)
	}()

	go func() {
		wg.Wait()
		close(jobsChan)
	}()

	select {
	case err := <-errChan:
		return nil, err
	case <-doneChan:
		return jobs, nil
	}
}

func processRows(rows <-chan []string, jobs chan<- Job, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for row := range rows {
		job := Job{
			ExternalID:  row[0],
			Title:       row[1],
			Description: row[2],
			Company:     row[3],
			Location:    row[4],
			URL:         row[5],
			Source:      3,
			Type:        "full-time",
		}

		if len(row) > 7 {
			job.Type = row[7]
		}

		jobs <- job
	}
}
