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

func InsertNotebook(ctx context.Context, notebook model.Notebook, collection dbInterface.CollectionAPI) (model.Notebook, *echo.HTTPError) {

	notebook.ID = primitive.NewObjectID()
	notebook.CreatedAt = time.Now()
	notebook.UpdatedAt = time.Now()
	_, err := collection.InsertOne(ctx, notebook)
	if err != nil {
		log.Errorf("Unable to insert the notebook :%+v", err)
		return notebook, echo.NewHTTPError(http.StatusInternalServerError, model.ErrorMessage{Message: "Unable to create the notebook"})
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
		return echo.NewHTTPError(http.StatusUnprocessableEntity, model.ErrorMessage{Message: "unable to find the notebook"})
	}
	updatedNotebook.Title = givenNotebook.Title
	//update the product, if err return 500
	_, err := collection.UpdateOne(ctx, filter, bson.M{"$set": updatedNotebook})
	if err != nil {
		log.Errorf("Unable to update the notebook : %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrorMessage{Message: "unable to update the notebook"})
	}
	return nil
}

// utils.ParseToken(c.Request().Header.Get("x-auth-token"))

func FindNotebooks(ctx context.Context, collection dbInterface.CollectionAPI,
	userId primitive.ObjectID) ([]model.Notebook, *echo.HTTPError) {

	var notebooks []model.Notebook

	cursor, err := collection.Find(ctx, bson.M{"user_id": userId})
	if err != nil {
		log.Errorf("Unable to find the notebooks : %v", err)
		return notebooks,
			echo.NewHTTPError(http.StatusNotFound, model.ErrorMessage{Message: "unable to find the notebooks"})
	}

	err = cursor.All(ctx, &notebooks)
	if err != nil {
		log.Errorf("Unable to read the cursor : %v", err)
		return notebooks,
			echo.NewHTTPError(http.StatusUnprocessableEntity, model.ErrorMessage{Message: "unable to parse retrieved notebooks"})
	}
	return notebooks, nil
}
