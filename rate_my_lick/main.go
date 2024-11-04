package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/topolovac/learning_projects/rate_my_lick/services"
)

type application struct {
	sampleService *services.SampleService
}

func main() {
	app := &application{
		sampleService: &services.SampleService{},
	}

	app.sampleService.CreateSample("Cool Song", "Description of very cool song.", "tintuntun_20241027213536")
	app.sampleService.CreateSample("Also Cool Song", "Artist X", "tintuntun_20241027213536")
	app.sampleService.CreateSample("Nothing Else Matters", "Metallica. Acustic version", "tintuntun_20241027213536")
	app.sampleService.CreateSample("Sip", "Tananana", "tintuntun_20241027213536")
	app.sampleService.CreateSample("Society Eddie Vedder", "", "tintuntun_20241027213536")

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(CreateSession)
	e.Static("/static", "static")

	e.GET("/", app.HomeHandler)
	e.GET("/create-lick", app.CreateLickHandler)
	e.POST("/publish-sample", app.PublishSampleHandler)

	e.POST("/lick/:id/rate/:rate", app.RateLickHandler)

	err := e.Start(":3000")
	e.Logger.Fatal(err)
}
