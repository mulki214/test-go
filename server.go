package main

import (
	"net/http"
	"io/ioutil"
	"net/url"
	"encoding/json"
   	"log"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
)
type (
	Matches struct {
		Competitionon	string 	`json:"competition"`
		Team1 			string 	`json:"team1"`
		Team2 			int 	`json:"team2"`
		Year 			int 	`json:"year"`
		Team1Goals		string	`json:"team1goals"`
		Team2Goals		string	`json:"team2goals"`
	}
)
type (
	Competion struct {
		Name 		string 	`json:"name"`
		Country 	string 	`json:"country"`
		Year 		int 	`json:"year"`
		Winner 		string 	`json:"winner"`
		Runnerup	string	`json:"runnerup"`
	}
)
type (
	Data struct {
		Page 		int 		`json:"page"`
		PerPage 	int 		`json:"per_page"`
		Total 		int 		`json:"total"`
		TotalPages 	int 		`json:"total_pages"`
		Data		[]Competion `json:"data"`
	}
)
type (
	DataMatches struct {
		Page 		int 		`json:"page"`
		PerPage 	int 		`json:"per_page"`
		Total 		int 		`json:"total"`
		TotalPages 	int 		`json:"total_pages"`
		Data		[]Matches 	`json:"data"`
	}
)

type (
	Input1 struct {
		Team	string 	`json:"team" validate:"required"` 
		Year 	string 	`json:"year" validate:"required"`
		Page 	string 	`json:"page" validate:"required"`
	}
)
type (
	Input2 struct {
		Name	string 	`json:"name" validate:"required"` 
		Page 	string 	`json:"page" validate:"required"`
	}
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

type CustomValidator struct {
    validator *validator.Validate
}
func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.validator.Struct(i)
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Throws unauthorized error
	if username != "admin" || password != "admin" {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		"Admin",
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func getCompetion(c echo.Context) error {
	name := c.QueryParam("competition")
	page := c.QueryParam("page")

	// Validator
	inpt := new(Input2)
	inpt.Name = name
	inpt.Page = page

	if err := c.Validate(inpt); err != nil {
        return err
    }

	Data := Data{}

	// Build URL
	baseUrl, err := url.Parse("https://jsonmock.hackerrank.com/api/football_competitions")
	params := url.Values{}
	params.Add("name", name)
	params.Add("page", page)

	// Add Query Parameters to the URL
	baseUrl.RawQuery = params.Encode()

	// Request to URL
	fmt.Println(baseUrl)
	resp, err := http.Get(baseUrl.String())
	if err != nil {
		log.Fatalln(err)
	}
  
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(body, &Data)
	return c.JSON(http.StatusOK, Data)
}

func getSchedule(c echo.Context) error {
	team := c.QueryParam("team")
	year := c.QueryParam("year")
	page := c.QueryParam("page")

	// Validator
	inpt := new(Input1)
	inpt.Team = team
	inpt.Year = year
	inpt.Page = page

	if err := c.Validate(inpt); err != nil {
        return err
    }

	Data := DataMatches{}
	// Data2 := DataMatches{}

	// Build URL
	baseUrl, err := url.Parse("https://jsonmock.hackerrank.com/api/football_matches")
	params := url.Values{}
	params.Add("team1", team)
	params.Add("year", year)
	params.Add("page", page)
	// Add Query Parameters to the URL
	baseUrl.RawQuery = params.Encode()
	// Request to URL
	fmt.Println(baseUrl)
	resp, err := http.Get(baseUrl.String())
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(body, &Data)

	return c.JSON(http.StatusOK, Data)
}



func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	r := e.Group("/api")

	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	r.Use(middleware.JWTWithConfig(config))

	// Routes
	r.GET("/competition", getCompetion)
	r.GET("/schedule", getSchedule)

	e.POST("/login", login)

	// Start server
	e.Logger.Fatal(e.Start(":8081"))
}