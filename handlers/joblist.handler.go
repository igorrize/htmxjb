package handlers

import (
	"github.com/igorrize/htmxjb/services"
	"github.com/igorrize/htmxjb/views/job_views"
	"github.com/labstack/echo/v4"

	"github.com/a-h/templ"
	"net/http"
)

type JobService interface {
	GetAllJobs() ([]services.Job, error)
}

type JobHandler struct {
	JobService JobService
}

func NewJobHandler(js JobService) *JobHandler {
	return &JobHandler{
		JobService: js,
	}
}

func (jh *JobHandler) jobListHandler(c echo.Context) error {
	c.Set("ISERROR", false)

	jobs, err := jh.JobService.GetAllJobs()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	titlePage := "Jobs List"
	return renderView(c, job_views.JobIndex(
		titlePage,
		job_views.JobList(titlePage, jobs),
	))
}

func renderView(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
