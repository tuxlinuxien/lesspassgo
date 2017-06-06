package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func registerPost(c echo.Context) error {
	var i struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	c.Bind(&i)
	return c.String(http.StatusOK, "")
}

func authPost(c echo.Context) error {
	var i struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	c.Bind(&i)
	return c.String(http.StatusOK, "")
}

func passwordsPost(c echo.Context) error {
	var i struct {
		Login     string `json:"login"`
		Site      string `json:"site"`
		Uppercase bool   `json:"uppercase"`
		Symbols   bool   `json:"symbols"`
		Lowercase bool   `json:"lowercase"`
		Numbers   bool   `json:"numbers"`
		Counter   int    `json:"counter"`
		Version   int    `json:"version"`
		Length    int    `json:"length"`
	}
	c.Bind(&i)
	return c.String(http.StatusOK, "")
}

func show(c echo.Context) error {
	var i interface{}
	c.Bind(&i)
	fmt.Println(i)
	return c.String(http.StatusOK, "")
}

// Start .
func Start(dbPath string, port int) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.POST("/api/auth/register/", registerPost)
	e.POST("/api/tokens/auth/", authPost)
	e.POST("/api/passwords/", passwordsPost)
	e.GET("/api/passwords/", show)
	e.Start(":" + strconv.Itoa(port))
}
