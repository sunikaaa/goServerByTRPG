package handler

// SayHello ...
// this is comment
import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"trpg.com/key"
	"trpg.com/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type jwtCustomClaims struct {
	UID   int    `json:"uid"`
	Email string `json:"email"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

var signingKey = []byte("secret")

// Config ... jwt Config
var Config = middleware.JWTConfig{
	Claims:     &jwtCustomClaims{},
	SigningKey: signingKey,
}

// Signup ... user create signup
func Signup(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		fmt.Printf("not passed this")
		return err
	}

	if user.Name == "" || user.Password == "" || user.Email == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid name or password",
		}
	}

	if u := model.FindUser(&model.User{Email: user.Email}); u.ID != 0 {
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "email already exists",
		}
	}
	pass := user.Password
	user.Password, _ = Generate(pass)

	model.CreateUser(user)
	Us := model.FindUser(&model.User{Email: user.Email})
	t, _ := createJWTToken(Us)
	// user.Password = ""

	return c.JSON(http.StatusCreated, map[string]string{
		"email": user.Email,
		"id":    fmt.Sprint(user.ID),
		"token": t,
	})
}

// Login ... this is user login
func Login(c echo.Context) error {
	u := new(model.User)
	if err := c.Bind(u); err != nil {
		return err
	}

	user := model.FindUser(&model.User{Email: u.Email})
	if user.ID == 0 {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "invalid email",
		}
	}
	tf, err := Compare(user.Password, u.Password)
	if !tf || err != nil {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "invalid password",
		}
	}
	t, _ := createJWTToken(user)

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
		"email": user.Email,
		"ID":    fmt.Sprint(user.ID),
	})
}

//userIDFromToken ... this is id and token
func userIDFromToken(c echo.Context) int {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	uid := claims.UID
	return uid
}

//Restricted ... this is testingcode
func Restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["Email"].(string)
	fmt.Println("this is run yet")
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

// func ComparePassword(c string) error {
// 	password := []byte("password")

// 	hashed, _ := bcrypt.GenerateFromPassword(password, 10)

// 	err := bcrypt.CompareHashAndPassword(hashed, password)

// }

// Generate ... this is generate hashed password
func Generate(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Compare ... comparing password
func Compare(hash, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, err
		}
		return false, err
	}
	return true, nil
}

func loadingKey() {

}

const (
	authorizeEndpoint = "https://accounts.google.com/o/oauth2/v2/auth"
	tokenEndpoint     = "https://www.googleapis.com/oauth2/v4/token"
	gcpScope          = "https://www.googleapis.com/auth/cloud-platform"
)

// CallbackRequest コールバックリクエスト
type CallbackRequest struct {
	Code  string `form:"code" query:"code"`
	State string `form:"state" query:"state"`
}

func oauthConfig(url string) *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     key.GoogleID,     // from https://console.developers.google.com/project/<your-project-id>/apiui/credential
		ClientSecret: key.GoogleSecret, // from https://console.developers.google.com/project/<your-project-id>/apiui/credential
		Endpoint:     google.Endpoint,
		RedirectURL:  url,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	}
	return config
}

// SignupByGoogle ... this is signup by OAuth2
func SignupByGoogle(c echo.Context) error {
	config := oauthConfig("http://localhost:3433/oauth")

	url := config.AuthCodeURL("state")
	fmt.Println(url)
	return c.JSON(http.StatusOK, map[string]string{
		"url": url,
	})
}

func SignUpByGoogleWithToken(c echo.Context) error {
	r := new(CallbackRequest)
	config := oauthConfig("http://localhost:3433/oauth")
	if err := c.Bind(r); err != nil {
		return err
	}
	if r.State != "state" {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "invalid password",
		}
	}
	fmt.Println(r)
	ctx := context.Background()
	tok, err := config.Exchange(ctx, r.Code)
	if err != nil {
		return err
	}
	client := config.Client(ctx, tok)
	resp, err := client.Get("https://oauth2.googleapis.com/tokeninfo")
	if err != nil {
		log.Fatalf("client get error")
	}
	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteArray))
	user := new(model.User)
	if err := json.Unmarshal(byteArray, &user); err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)
	if u := model.FindUser(&model.User{Email: user.Email}); u.ID != 0 {
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "email already exists",
		}
	}

	t, _ := createJWTToken(*user)
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func createJWTToken(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Email": user.Email,
		"ID":    user.ID,
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	t, err := token.SignedString([]byte(key.Secret))
	if err != nil {
		log.Fatalln(err)
	}
	return t, err
}
