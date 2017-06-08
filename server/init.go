package server

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func checkAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("authorization")
		if token == "" {
			return c.String(http.StatusUnauthorized, "")
		}
		if !strings.HasPrefix(token, "JWT ") {
			return c.String(http.StatusUnauthorized, "")
		}
		token = strings.Replace(token, "JWT ", "", 1)
		user, err := checkToken(token)
		if err != nil {
			return c.String(http.StatusUnauthorized, err.Error())
		}
		c.Set("user", user)
		return next(c)
	}
}

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
	userInterface := c.Get("user")
	user := userInterface.(*UserModel)
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
	if err := c.Bind(&i); err != nil {
		log.Println(err)
		return err
	}
	p, err := CreatePassword(user.ID, i.Login, i.Site, i.Uppercase, i.Lowercase, i.Symbols, i.Numbers, i.Counter, i.Version, i.Length)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, p)
}

func show(c echo.Context) error {
	var i interface{}
	c.Bind(&i)
	return c.String(http.StatusOK, "")
}

func passwordsGet(c echo.Context) error {
	userInterface := c.Get("user")
	user := userInterface.(*UserModel)
	var result struct {
		Count   int             `json:"count"`
		Next    interface{}     `json:"next"`
		Prev    interface{}     `json:"prev"`
		Results []PasswordModel `json:"results"`
	}
	passwords := GetPasswordsByUserID(user.ID)
	result.Count = len(passwords)
	result.Results = passwords
	return c.JSON(http.StatusOK, &result)
}

func passwordsDelete(c echo.Context) error {
	userInterface := c.Get("user")
	user := userInterface.(*UserModel)
	passwordID := c.Param("id")
	err := DeletePasswordByIDAndUserID(passwordID, user.ID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": http.StatusText(http.StatusUnauthorized),
		})
	}
	return c.String(http.StatusOK, "")
}

func passwordsPut(c echo.Context) error {
	userInterface := c.Get("user")
	user := userInterface.(*UserModel)
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
	if err := c.Bind(&i); err != nil {
		log.Println(err)
		return err
	}
	passwordID := c.Param("id")
	p, err := GetPasswordByID(passwordID)
	if err != nil {
		log.Println(err)
		return nil
	}
	if p.UserID != user.ID {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": http.StatusText(http.StatusUnauthorized),
		})
	}
	p.Login = i.Login
	p.Site = i.Site
	p.Uppercase = i.Uppercase
	p.Symbols = i.Symbols
	p.Lowercase = i.Lowercase
	p.Numbers = i.Numbers
	p.Counter = i.Counter
	p.Version = i.Version
	p.Length = i.Length
	if err := p.Update(); err != nil {
		log.Println(err)
		return err
	}
	return c.JSON(http.StatusOK, p)
}

// Start .
func Start(dbPath string, port int) {
	openDB(dbPath)
	createTable()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.AddTrailingSlash())
	e.POST("/api/auth/register/", registerPost)
	e.POST("/api/tokens/auth/", authPost)
	e.POST("/api/passwords/", passwordsPost, checkAuth)
	e.GET("/api/passwords/", passwordsGet, checkAuth)
	e.PUT("/api/passwords/:id/", passwordsPut, checkAuth)
	e.DELETE("/api/passwords/:id/", passwordsDelete, checkAuth)
	e.Start(":" + strconv.Itoa(port))
}
