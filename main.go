package main

import (
	"context"
	"fmt"

	"github.com/go-study-projs/vue-evernote-api/handler"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-study-projs/vue-evernote-api/config"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/labstack/gommon/random"
)

var (
	db       *mongo.Database
	usersCol *mongo.Collection
	cfg      config.Properties
)

const (
	//CorrelationID is a request id unique to the request being made
	CorrelationID = "X-Correlation-ID"
)

func init() {
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Configuration cannot be read : %v", err)
	}

	ctx := context.Background()
	connectURI := fmt.Sprintf("mongodb://%s:%s", cfg.DBHost, cfg.DBPort)
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(connectURI))
	if err != nil {
		log.Fatalf("Unable to connect to database : %v", err)
	}
	db = c.Database(cfg.DBName)
	usersCol = db.Collection(cfg.UserCollection)

	// add db index to username
	isUserIndexUnique := true
	indexModel := mongo.IndexModel{
		Keys: bson.M{"username": 1},
		Options: &options.IndexOptions{
			Unique: &isUserIndexUnique,
		},
	}
	_, err = usersCol.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Fatalf("Unable to create an index : %+v", err)
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
	uh := &handler.UsersHandler{Col: usersCol}
	e.POST("/auth/register", uh.CreateUser)
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
