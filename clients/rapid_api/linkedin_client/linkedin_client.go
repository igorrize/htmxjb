package linkedin_client

import (
	"fmt"
	"io"
	"net/http"
)

type LinkedinClient interface {
	GetJobs() error
}

type LinkedinClient struct {
	apiKey string
}

func NewLinkedinClient(apiKey string) *LinkedinClient {
	return &LinkedinClient{
		apiKey: apiKey,
	}
}
func(c *LinkedinClient) GetJobs() error {

	url := "https://linkedin-job-search-api.p.rapidapi.com/active-jb-7d"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-key", "3ba6283321msh80a967930e01301p13ea9fjsna2afafa73bdc")
	req.Header.Add("x-rapidapi-host", "linkedin-job-search-api.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}
