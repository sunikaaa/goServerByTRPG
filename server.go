package main

import (
	"crypto/subtle"
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	db, err := sql.Open("mysql", "root:sunica@/trpg")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// e.Use(middleware.JWTWithConfig(handler.Config))
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
			return true, nil
		}
		return false, nil
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// e.POST("/login", createUser)
	// e.POST("/signup", createUser)
	e.Logger.Fatal(e.Start(":1323"))
}

// type User [
// 	username string
// 	password string
// 	email string
// ]

// func createUser(c echo.Context) error {
// 	var form User
// 	if err := c.
// 	return c.String(http.StatusOK, "name:"+name+", email:"+email)
// }
