package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

func main() {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.Static("/static", "assets")

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
	m := &MessageError{
		Message: "Match not found",
	}
	return c.JSON(http.StatusNotFound, m)
}

type Match struct {
	HomeTeam     string `json:"homeTeam,omitempty"`
	AwayTeam     string `json:"awayTeam,omitempty"`
	Championship string `json:"championship,omitempty"`
}

type MessageError struct {
	Message string `json:"message,omitempty"`
}
