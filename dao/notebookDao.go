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
