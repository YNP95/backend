package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"ynp/api"
	"ynp/env"

	"github.com/golang-jwt/jwt"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "ynp/docs"
)

type Stats struct {
	Uptime  time.Time     `json:"uptime"`
	ExpTime time.Duration `json:"exptime"`
}

func NewStats() *Stats {
	return &Stats{
		Uptime: time.Now(),
	}
}

func (s *Stats) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)

		s.ExpTime = time.Since(s.Uptime)
		log.Println(c.Path(), s.ExpTime)
		return err
	}
}

func init() {
	env.MyDB = api.NewDb()
}

// @title YNP SERVER API
// @version 1.0
// @BasePath /v1
func main() {
	defer api.CloseDb(env.MyDB)

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Use(echojwt.JWT([]byte("secret")))

	e.GET("/v1/jwttest", func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token) // by default token is stored under `user` key
		if !ok {
			return errors.New("JWT token missing or invalid")
		}
		claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
		if !ok {
			return errors.New("failed to cast claims as jwt.MapClaims")
		}
		return c.JSON(http.StatusOK, claims)
	})
	// e.Use(echojwt.WithConfig(echojwt.Config{
	// 	// ...
	// 	SigningKey:             []byte("secret"),
	// 	// ...
	//   }))

	e.Use(middleware.CORS())

	e.GET("/v1/", api.Index)
	e.GET("/v1/swagger/*", echoSwagger.WrapHandler)

	e.GET("/v1/random", api.Random)

	e.GET("/v1/users/get/:name", api.GetUserInfo)
	e.GET("/v1/users/exist/:name", api.IdDuplicateCheck)
	e.POST("/v1/users/signup", api.NewUserInfo)
	e.POST("/v1/users/signin", api.SignIn)

	e.POST("/v1/table/create", api.CreateTable)

	e.GET("/v1/crawl/lotto/:round", api.CrawlingLottoNum)
	e.GET("/v1/crawl/lotto/all", api.CrawlingLottoNumAll)

	e.GET("/v1/lotto/get/:round", api.GetLottoNum)

	e.Logger.Fatal(e.Start(":8080"))
}
