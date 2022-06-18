package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/rs/zerolog"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var log *zerolog.Logger

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Caller().Logger()
	log = &logger
}

func main() {
	e := echo.New()
	// Middleware
	e.Logger.SetOutput(ioutil.Discard)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			log.Debug().
				Interface("headers", req.Header).
				Msg(">>> " + req.Method + " " + req.RequestURI)
			if err = next(c); err != nil {
				c.Error(err)
			}
			log.Debug().
				Str("latency", time.Now().Sub(start).String()).
				Int("status", res.Status).
				Interface("headers", res.Header()).
				Msg("<<< " + req.Method + " " + req.RequestURI)
			return
		}
	})
	e.Use(middleware.Recover())
	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.Static("/static", "assets/api-docs")

	// Server
	e.GET("/api/matches/:id", GetMatch)
	e.GET("/health", Health)
	e.Logger.Fatal(e.Start(":9999"))
}

func Health(c echo.Context) error {
	return c.JSON(200, &HealthData{Status: "UP"})
}

type HealthData struct {
	Status string `json:"status,omitempty"`
}

func GetMatch(c echo.Context) error {
	m := &Match{
		HomeTeam:     "Barcelona",
		AwayTeam:     "Real Madrid",
		Championship: "UEFA",
	}
	return c.JSON(http.StatusOK, m)
}

type Match struct {
	HomeTeam     string `json:"homeTeam,omitempty"`
	AwayTeam     string `json:"awayTeam,omitempty"`
	Championship string `json:"championship,omitempty"`
}
