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
	db          *mongo.Database
	userCol     *mongo.Collection
	notebookCol *mongo.Collection
	noteCol     *mongo.Collection
	cfg         config.Properties
)

const (
	//CorrelationID is a request id unique to the request being made
	CorrelationID          = "X-Correlation-ID"
	UserCollectionName     = "user"
	NotebookCollectionName = "notebook"
	NoteCollectionName     = "note"
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

	userCol = db.Collection(UserCollectionName)
	notebookCol = db.Collection(NotebookCollectionName)
	noteCol = db.Collection(NoteCollectionName)

	// add db index to username
	isUserIndexUnique := true
	indexModel := mongo.IndexModel{
		Keys: bson.M{"username": 1},
		Options: &options.IndexOptions{
			Unique: &isUserIndexUnique,
		},
	}
	_, err = userCol.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Fatalf("Unable to create an index : %+v", err)
	}
}

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	e.Pre(middleware.RemoveTrailingSlash())
	e.Pre(addCorrelationID)
	jwtMiddleware := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(cfg.JwtTokenSecret),
	})
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${time_rfc3339_nano} ${remote_ip} ${header:X-Correlation-ID} ${host} ${method} ${uri} ${user_agent} ` +
			`${status} ${error} ${latency_human}` + "\n",
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{"*"},
		MaxAge:       3600, //Access-Control-Max-Age,3600s内客户端只会发送一次pre-flight请求
	}))

	uh := &handler.UserHandler{Col: userCol}
	e.POST("/auth/register", uh.CreateUser)
	e.POST("/auth/login", uh.Login)

	nbh := &handler.NotebookHandler{Col: notebookCol, NCol: noteCol}
	e.GET("/notebooks", nbh.GetNotebooks, jwtMiddleware)
	e.POST("/notebooks", nbh.CreateNotebook, middleware.BodyLimit("1M"), jwtMiddleware)
	e.PATCH("/notebooks/:notebookId", nbh.UpdateNoteBook, middleware.BodyLimit("1M"), jwtMiddleware)
	e.DELETE("/notebooks/:notebookId", nbh.DeleteNoteBook, jwtMiddleware)

	nh := &handler.NoteHandler{Col: noteCol}
	e.POST("/notes/to/:notebookId", nh.CreateNote, middleware.BodyLimit("4M"), jwtMiddleware)
	e.GET("/notes/from/:notebookId", nh.GetNotes, jwtMiddleware)
	e.DELETE("/notes/:noteId", nh.MoveToTrash, jwtMiddleware)
	e.PATCH("/notes/:noteId", nh.UpdateNote, middleware.BodyLimit("4M"), jwtMiddleware)
	e.DELETE("/notes/:noteId/confirm", nh.DeleteNote, jwtMiddleware)
	e.PATCH("/notes/:noteId/revert", nh.RevertNote, jwtMiddleware)
	e.GET("/notes/trash", nh.GetNotesInTrash, jwtMiddleware)

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
