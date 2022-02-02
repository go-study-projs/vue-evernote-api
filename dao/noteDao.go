package dao

import (
	"context"
	"net/http"
	"time"

	"github.com/go-study-projs/vue-evernote-api/dbInterface"
	"github.com/go-study-projs/vue-evernote-api/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertNote(ctx context.Context, note model.Note, collection dbInterface.CollectionAPI) (model.Note, *echo.HTTPError) {

	note.ID = primitive.NewObjectID()
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()
	_, err := collection.InsertOne(ctx, note)
	if err != nil {
		log.Errorf("Unable to insert the note :%+v", err)
		return note, echo.NewHTTPError(http.StatusInternalServerError, model.ErrorMessage{Message: "Unable to create the note"})
	}
	return note, nil
}
