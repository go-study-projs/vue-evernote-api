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

type NoteOperateType int

const (
	SoftDelete NoteOperateType = iota
	Revert
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

func FindNotesInTrash(ctx context.Context, collection dbInterface.CollectionAPI,
	userId primitive.ObjectID) ([]model.Note, *echo.HTTPError) {

	var notes []model.Note

	cursor, err := collection.Find(ctx, bson.M{"user_id": userId, "is_deleted": true})
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

func ModifyNote(ctx context.Context, collection dbInterface.CollectionAPI,
	noteId primitive.ObjectID, givenNote model.Note) *echo.HTTPError {

	var updatedNote model.Note
	filter := bson.M{"_id": noteId}

	res := collection.FindOne(ctx, bson.M{"_id": noteId})
	if err := res.Decode(&updatedNote); err != nil {
		log.Errorf("unable to decode to note :%v", err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, model.ErrorMessage{Message: "unable to find the note"})
	}

	updatedNote.Title = givenNote.Title
	updatedNote.Content = givenNote.Content

	_, err := collection.UpdateOne(ctx, filter, bson.M{"$set": updatedNote})
	if err != nil {
		log.Errorf("Unable to update the note : %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrorMessage{Message: "unable to update the note"})
	}
	return nil
}

func SoftDeleteOrRevertNote(ctx context.Context, collection dbInterface.CollectionAPI,
	noteId primitive.ObjectID, operateType NoteOperateType) *echo.HTTPError {

	var note model.Note
	filter := bson.M{"_id": noteId}

	res := collection.FindOne(ctx, bson.M{"_id": noteId})
	if err := res.Decode(&note); err != nil {
		log.Errorf("unable to decode to note :%v", err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, model.ErrorMessage{Message: "unable to find the note"})
	}

	switch operateType {
	case SoftDelete:
		note.IsDeleted = true
	case Revert:
		note.IsDeleted = false
	}

	_, err := collection.UpdateOne(ctx, filter, bson.M{"$set": note})
	if err != nil {
		log.Errorf("Unable to update the note : %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrorMessage{Message: "unable to move the note to trash"})
	}
	return nil
}

func DeleteNote(ctx context.Context, collection dbInterface.CollectionAPI, noteId primitive.ObjectID) (int64, *echo.HTTPError) {

	res, err := collection.DeleteOne(ctx, bson.M{"_id": noteId})
	if err != nil {
		log.Errorf("Unable to delete the note : %v", err)
		return 0,
			echo.NewHTTPError(http.StatusInternalServerError, model.ErrorMessage{Message: "unable to delete the note"})
	}
	return res.DeletedCount, nil
}
