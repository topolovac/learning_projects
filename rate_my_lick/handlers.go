package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/topolovac/learning_projects/rate_my_lick/components"
)

var GuestUserId = "guest_user_id"

func CreateSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := c.Cookie(GuestUserId)
		if err != nil {
			cookie := new(http.Cookie)
			cookie.Name = GuestUserId
			cookie.Value = uuid.New().String()
			c.SetCookie(cookie)
		}
		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
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
	return Render(c, http.StatusOK, components.Home(app.sampleService.GetSamples()))
}

func (app *application) CreateLickHandler(c echo.Context) error {
	return Render(c, http.StatusOK, components.CreateLick())
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
	dst, err := os.Create("./static/licks/" + filename)
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

func (app *application) RateLickHandler(c echo.Context) error {
	guest_user_id, err := c.Cookie(GuestUserId)
	if err != nil {
		return err
	}

	lickId := c.Param("id")
	rate := c.Param("rate")

	r, err := strconv.Atoi(rate)
	if err != nil {
		return err
	}
	sample, err := app.sampleService.RateSample(uuid.MustParse(lickId), r, uuid.MustParse(guest_user_id.Value))
	if err != nil {
		return err
	}

	return Render(c, http.StatusOK, components.RatingLabel(len(sample.Ratings[r])))
}
