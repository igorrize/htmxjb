package responses

import "time"

type IndeedResponse struct {
    State       string `json:"state"`
    Name        string `json:"name"`
    Data        IndeedData `json:"data"`
    ID          string `json:"id"`
    Progress    int    `json:"progress"`
    ReturnValue IndeedReturnValue `json:"returnvalue"`
}

type IndeedData struct {
    Scraper IndeedScraperParams `json:"scraper"`
}

type IndeedScraperParams struct {
    Query     string `json:"query"`
    Location  string `json:"location"`
    JobType   string `json:"jobType"`
    Radius    string `json:"radius"`
    Sort      string `json:"sort"`
    FromDays  string `json:"fromDays"`
    Country   string `json:"country"`
}

type IndeedReturnValue struct {
    Data []IndeedJob `json:"data"`
}

type IndeedJob struct {
    Title           string    `json:"title"`
    JobType         string    `json:"jobType"`
    CompanyName     string    `json:"companyName"`
    CompanyUrl      string    `json:"companyUrl"`
    Location        IndeedLocation `json:"location"`
    Attributes      []string  `json:"attributes"`
    DescriptionHtml string    `json:"descriptionHtml"`
    DescriptionText string    `json:"descriptionText"`
    Age            string    `json:"age"`
    DatePublished  time.Time `json:"datePublished"`
    JobKey         string    `json:"jobKey"`
    JobUrl         string    `json:"jobUrl"`
    RemoteLocation bool      `json:"remoteLocation"`
}

type IndeedLocation struct {
    CountryCode          string  `json:"countryCode"`
    City                 string  `json:"city"`
    Latitude            float64 `json:"latitude"`
    Longitude           float64 `json:"longitude"`
    FormattedAddressLong  string  `json:"formattedAddressLong"`
    FormattedAddressShort string  `json:"formattedAddressShort"`
}
