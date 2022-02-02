package dao

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"

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

func FindNotes(ctx context.Context, collection dbInterface.CollectionAPI,
	notebookId primitive.ObjectID) ([]model.Note, *echo.HTTPError) {

	var notes []model.Note

	cursor, err := collection.Find(ctx, bson.M{"notebook_id": notebookId, "is_deleted": false})
	if err != nil {
		log.Errorf("Unable to find the notes : %v", err)
		return notes,
			echo.NewHTTPError(http.StatusNotFound, model.ErrorMessage{Message: "unable to find the notes"})
	}

	err = cursor.All(ctx, &notes)
	if err != nil {
		log.Errorf("Unable to read the cursor : %v", err)
		return notes,
			echo.NewHTTPError(http.StatusUnprocessableEntity, model.ErrorMessage{Message: "unable to parse retrieved notes"})
	}
	return notes, nil
}
