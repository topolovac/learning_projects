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
	userId, err := getUserIdFromCookie(c)
	if err != nil {
		return err
	}
	return Render(c, http.StatusOK, components.HomePage(app.sampleService.GetSamplesByRating(), userId))
}

func (app *application) LatestPageHandler(c echo.Context) error {
	userId, err := getUserIdFromCookie(c)
	if err != nil {
		return err
	}
	return Render(c, http.StatusOK, components.LatestPage(app.sampleService.GetSamplesOrderByLatest(), userId))
}

func (app *application) CreateLickHandler(c echo.Context) error {
	return Render(c, http.StatusOK, components.CreateLick())
}

func (app *application) PublishSampleHandler(c echo.Context) error {
	name := c.FormValue("songname")
	description := c.FormValue("songdescription")

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

	sample, err := app.sampleService.CreateSample(name, description, filename)
	if err != nil {
		return err
	}

	c.Response().Header().Add("HX-Redirect", fmt.Sprintf("/lick/%s", sample.Id.String()))
	return nil
}

func (app *application) RateLickHandler(c echo.Context) error {
	user_id, err := getUserIdFromCookie(c)
	if err != nil {
		return err
	}

	lickId := c.Param("id")
	rate := c.Param("rate")

	r, err := strconv.Atoi(rate)
	if err != nil {
		return err
	}

	sample, err := app.sampleService.RateSample(uuid.MustParse(lickId), r, user_id)
	if err != nil {
		return err
	}

	return Render(c, http.StatusOK, components.RatingSection(*sample, user_id))
}

func (app *application) LickHandler(c echo.Context) error {
	id := c.Param("id")
	lickId, err := uuid.Parse(id)
	if err != nil {
		return nil
	}

	sample, err := app.sampleService.GetSampleById(lickId)
	if err != nil {
		return err
	}

	userId, err := getUserIdFromCookie(c)
	if err != nil {
		return nil
	}

	return Render(c, http.StatusOK, components.LickPage(*sample, userId))
}

// UTILS
func getUserIdFromCookie(c echo.Context) (uuid.UUID, error) {
	cookie, err := c.Cookie(GuestUserId)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuid.Parse(cookie.Value)
}
