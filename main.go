package main

import (
	"fmt"
	"github.com/go-study-projs/vue-evernote-api/config"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/labstack/gommon/random"
)

var (
	cfg config.Properties
)

const (
	//CorrelationID is a request id unique to the request being made
	CorrelationID = "X-Correlation-ID"
)

func init() {
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Configuration cannot be read : %v", err)
	}
}

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	e.Pre(middleware.RemoveTrailingSlash())
	e.Pre(addCorrelationID)
	//jwtMiddleware := middleware.JWTWithConfig(middleware.JWTConfig{
	//	SigningKey:  []byte(cfg.JwtTokenSecret),
	//	TokenLookup: "header:x-auth-token",
	//})

	e.Logger.Infof("Listening on %s:%s", cfg.Host, cfg.Port)
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)))
}

func addCorrelationID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// generate correlation id
		id := c.Request().Header.Get(CorrelationID)
		var newID string
		if id == "" {
			//generate a random number
			newID = random.String(12)
		} else {
			newID = id
		}
		c.Request().Header.Set(CorrelationID, newID)
		c.Response().Header().Set(CorrelationID, newID)
		return next(c)
	}
}
