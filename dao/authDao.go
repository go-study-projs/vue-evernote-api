package dao

import (
	"context"
	"net/http"

	"github.com/go-study-projs/vue-evernote-api/dbInterface"
	"github.com/go-study-projs/vue-evernote-api/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func InsertUser(ctx context.Context, user model.User, collection dbInterface.CollectionAPI) (model.User, *echo.HTTPError) {
	var newUser model.User
	res := collection.FindOne(ctx, bson.M{"username": user.Username})
	err := res.Decode(&newUser)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Errorf("Unable to decode retrieved user: %v", err)
		return newUser,
			echo.NewHTTPError(http.StatusUnprocessableEntity, model.ErrorMessage{Message: "Unable to decode retrieved user"})
	}
	if newUser.Username != "" {
		log.Errorf("User  %s already exists", user.Username)
		return newUser,
			echo.NewHTTPError(http.StatusBadRequest, model.ErrorMessage{Message: "User already exists"})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		log.Errorf("Unable to hash the password: %v", err)
		return newUser,
			echo.NewHTTPError(http.StatusInternalServerError, model.ErrorMessage{Message: "Unable to process the password"})
	}
	user.Password = string(hashedPassword)
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		log.Errorf("Unable to insert the user :%+v", err)
		return newUser,
			echo.NewHTTPError(http.StatusInternalServerError, model.ErrorMessage{Message: "Unable to create the user"})
	}
	return model.User{Username: user.Username}, nil
}
