package dao

import (
	"context"
	"net/http"
	"time"

	"github.com/go-study-projs/vue-evernote-api/utils"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/go-study-projs/vue-evernote-api/dbInterface"
	"github.com/go-study-projs/vue-evernote-api/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertNotebook(ctx context.Context, notebook model.Notebook, collection dbInterface.CollectionAPI) (model.Notebook, *echo.HTTPError) {

	notebook.ID = primitive.NewObjectID()
	notebook.CreatedAt = time.Now()
	notebook.UpdatedAt = time.Now()
	_, err := collection.InsertOne(ctx, notebook)
	if err != nil {
		log.Errorf("Unable to insert the notebook :%+v", err)
		return notebook, echo.NewHTTPError(http.StatusInternalServerError, model.Response{Msg: "Unable to create the notebook"})
	}
	return notebook, nil
}

func ModifyNoteBook(ctx context.Context,
	collection dbInterface.CollectionAPI,
	notebookId primitive.ObjectID,
	givenNotebook model.Notebook) *echo.HTTPError {
	var updatedNotebook model.Notebook
	filter := bson.M{"_id": notebookId}
	res := collection.FindOne(ctx, filter)
	if err := res.Decode(&updatedNotebook); err != nil {
		log.Errorf("unable to decode to notebook :%v", err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, model.Response{Msg: "unable to find the notebook"})
	}
	updatedNotebook.Title = givenNotebook.Title
	_, err := collection.UpdateOne(ctx, filter, bson.M{"$set": updatedNotebook})
	if err != nil {
		log.Errorf("Unable to update the notebook : %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, model.Response{Msg: "unable to update the notebook"})
	}
	return nil
}

func FindNotebooks(ctx context.Context, collection dbInterface.CollectionAPI,
	userId primitive.ObjectID) ([]model.NotebookResponse, *echo.HTTPError) {

	var NotebookResponses []model.NotebookResponse

	pipeline := utils.NewAggregatePipe().
		Match(bson.M{"user_id": userId}).
		LookupOne("note", "notes_info", "_id", "notebook_id").
		Project(bson.M{
			"title":       "$title",
			"created_at":  "$created_at",
			"updated_at":  "$updated_at",
			"user_id":     "$user_id",
			"note_counts": bson.M{"$size": "$notes_info"},
		}).
		QueryM()

	cursor, err := collection.Aggregate(ctx, pipeline)

	if err != nil {
		log.Errorf("Unable to find the notebooks : %v", err)
		return NotebookResponses,
			echo.NewHTTPError(http.StatusNotFound, model.Response{Msg: "unable to find the notebooks"})
	}

	err = cursor.All(ctx, &NotebookResponses)
	if err != nil {
		log.Errorf("Unable to read the cursor : %v", err)
		return NotebookResponses,
			echo.NewHTTPError(http.StatusUnprocessableEntity, model.Response{Msg: "unable to parse retrieved notebooks"})
	}

	if NotebookResponses == nil {
		NotebookResponses = []model.NotebookResponse{}
	}

	return NotebookResponses, nil
}

func DeleteNoteBook(ctx context.Context, collection dbInterface.CollectionAPI,
	noteCollection dbInterface.CollectionAPI, notebookId primitive.ObjectID) (int64, *echo.HTTPError) {

	var notes []model.Note

	cursor, err := noteCollection.Find(ctx, bson.M{"notebook_id": notebookId})
	if err != nil {
		log.Errorf("Unable to find the notes : %v", err)
		return 0, echo.NewHTTPError(http.StatusNotFound, model.Response{Msg: "unable to find the notes"})
	}

	err = cursor.All(ctx, &notes)
	if err != nil {
		log.Errorf("Unable to read the cursor : %v", err)
		return 0, echo.NewHTTPError(http.StatusUnprocessableEntity, model.Response{Msg: "unable to parse retrieved notes"})
	}

	if notes != nil && len(notes) > 0 {
		return 0, echo.NewHTTPError(http.StatusBadRequest, model.Response{Msg: "????????????????????????????????????????????????????????????????????????"})
	}

	res, err := collection.DeleteOne(ctx, bson.M{"_id": notebookId})
	if err != nil {
		log.Errorf("Unable to delete the notebook : %v", err)
		return 0,
			echo.NewHTTPError(http.StatusInternalServerError, model.Response{Msg: "unable to delete the notebook"})
	}
	return res.DeletedCount, nil
}
