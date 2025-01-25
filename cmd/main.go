package main

import (
	"github.com/igorrize/htmxjb/db"
	"github.com/igorrize/htmxjb/handlers"
	"github.com/igorrize/htmxjb/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	DB_NAME string = "jobs.db"
)

func main() {

	e := echo.New()

	e.Static("/tailwind/public", "./tailwind/public")

	e.Static("/assets", "./assets")
	// Helpers Middleware
	e.Use(middleware.Logger())

	store, err := db.NewStore(DB_NAME)
	if err != nil {
		e.Logger.Fatal(err)
	}

	js := services.NewJobServices(services.Job{}, store)
	jh := handlers.NewJobHandler(js)
	// Setting Routes
	handlers.SetupRoutes(e, jh)

	// Start Server
	e.Logger.Fatal(e.Start(":8080"))
}
