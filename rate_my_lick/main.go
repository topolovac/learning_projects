package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	"github.com/topolovac/learning_projects/rate_my_lick/components"
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
	e.Static("/static", "static")

	e.GET("/", app.HomeHandler)
	e.GET("licks", app.LicksHandler)
	e.POST("/publish-sample", app.PublishSampleHandler)

	err := e.Start(":3000")
	e.Logger.Fatal(err)
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func (app *application) HomeHandler(c echo.Context) error {
	return Render(c, http.StatusOK, components.Home())
}

func (app *application) LicksHandler(c echo.Context) error {
	return Render(c, http.StatusOK, components.Licks(app.sampleService.GetSamples()))
}

func (app *application) PublishSampleHandler(c echo.Context) error {
	name := c.FormValue("songname")
	description := c.FormValue("songdescription")
	fmt.Println("name=" + name)
	fmt.Println("description=" + description)

	file, err := c.FormFile("audiosample")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	timestamp := time.Now()
	filename := name + "_" + timestamp.Format("20060102150405")
	dst, err := os.Create("./static/" + filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	app.sampleService.CreateSample(name, description, filename)

	return Render(c, http.StatusOK, components.PublishSample())
}
