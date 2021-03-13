package main

import (
	"fmt"
	"net/http"

	"trpg.com/handler"
	"trpg.com/key"
	"trpg.com/model"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//Test ... this
type Test struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db := model.InitDB()
	println(&db)
	println("this is ")
	model.TestDB(db)
	model.ShowAll()
	defer db.Close()
	// e.Use(middleware.JWTWithConfig(handler.Config))
	// e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	// 	if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
	// 		subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
	// 		return true, nil
	// 	}
	// 	return false, nil
	// }))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/test", display)

	e.POST("/signup", handler.Signup)

	e.POST("/signin", handler.Login)

	e.POST("/signupByGoogle", handler.SignupByGoogle)
	e.POST("/signUpByGoogleWithToken", handler.SignUpByGoogleWithToken)

	r := e.Group("/user")
	r.Use(middleware.JWT([]byte(key.Secret)))
	r.GET("", handler.Restricted)

	fmt.Println("http://localhost:1323")
	e.Logger.Fatal(e.Start(":1323"))
}
func display(c echo.Context) error {
	u := new(Test)
	if err := c.Bind(u); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, u)
}
