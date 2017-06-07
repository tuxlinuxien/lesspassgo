package server

import (
	"fmt"
	"log"
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
	if err := c.Bind(&i); err != nil {
		log.Println(err)
		return err
	}
	if err := CreateUser(i.Email, i.Password); err != nil {
		log.Println(err)
		return err
	}
	return c.String(http.StatusOK, "")
}

func authPost(c echo.Context) error {
	var i struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&i); err != nil {
		log.Println(err)
		return err
	}
	user, err := AuthUser(i.Email, i.Password)
	if err != nil {
		log.Println(err)
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": signUser(user.ID),
	})
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
	token := c.Request().Header.Get("authorization")
	fmt.Println("token:", token)
	var i interface{}
	c.Bind(&i)
	fmt.Println(i)
	return c.String(http.StatusOK, "")
}

// Start .
func Start(dbPath string, port int) {
	openDB(dbPath)
	createTable()
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
