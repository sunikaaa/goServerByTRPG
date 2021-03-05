package main

import (
	"crypto/subtle"
	"net/http"

	"trpg.com/handler"
	"trpg.com/model"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db := model.InitDB()

	e.Use(middleware.JWTWithConfig(handler.Config))
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
			return true, nil
		}
		return false, nil
	}))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	r := e.Group("/user")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", handler.Restricted)

	e.Logger.Fatal(e.Start(":1323"))
}
