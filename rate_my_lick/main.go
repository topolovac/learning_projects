package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/topolovac/learning_projects/rate_my_lick/services"
)

type application struct {
	sampleService *services.SampleService
}

func main() {
	app := &application{
		sampleService: &services.SampleService{},
	}

	e := echo.New()
	e.HTTPErrorHandler = app.CustomHTTPErrorHandler

	e.Use(middleware.Logger())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	customCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "custom_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
	)
	if err := prometheus.Register(customCounter); err != nil {
		log.Fatal(err)
	}

	e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		AfterNext: func(c echo.Context, err error) {
			customCounter.Inc()
		},
	}))

	e.Static("/static", "static")

	e.Use(CreateSession)

	e.GET("/", app.HomeHandler)
	e.GET("/create-lick", app.CreateLickHandler)
	e.GET("/latest", app.LatestPageHandler)
	e.POST("/publish-sample", app.PublishSampleHandler)
	e.GET("/lick/:id", app.LickHandler)
	e.GET("/my-licks", app.MyLicksPageHandler)

	e.POST("/lick/:id/rate/:rate", app.RateLickHandler)

	e.GET("/metrics", echoprometheus.NewHandler())

	err := e.Start(":3000")
	e.Logger.Fatal(err)
}
