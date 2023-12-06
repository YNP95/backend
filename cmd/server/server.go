package main

import (
	"log"
	"time"

	"ynp/api"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func main() {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	// s := NewStats()
	// e.Use(s.Middleware)
	e.Use(middleware.CORS())

	e.GET("/", api.Index)
	e.GET("/random", api.Random)

	e.GET("/users/get/:name", api.GetUserInfo)
	e.POST("/users/signup", api.NewUserInfo)
	e.POST("/users/signin", api.SignIn)

	e.POST("/table/create", api.CreateTable)

	e.GET("/get/lotto/:round", api.GetLottoNum)

	// e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
	// 	SigningKey:  []byte("secret"),
	// 	TokenLookup: "query:token",
	// }))

	e.Logger.Fatal(e.Start(":8080"))
}
