package dao

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-study-projs/vue-evernote-api/dbInterface"
	"github.com/go-study-projs/vue-evernote-api/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func InsertUser(ctx context.Context, collection dbInterface.CollectionAPI, user model.User) (model.User, *echo.HTTPError) {
	var existedUser model.User
	res := collection.FindOne(ctx, bson.M{"username": user.Username})
	err := res.Decode(&existedUser)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Errorf("Unable to decode retrieved user: %v", err)
		return existedUser,
			echo.NewHTTPError(http.StatusUnprocessableEntity, model.Response{Msg: "Unable to decode retrieved user"})
	}
	if existedUser.Username != "" {
		log.Errorf("User  %s already exists", user.Username)
		return existedUser,
			echo.NewHTTPError(http.StatusBadRequest, model.Response{Msg: "User already exists"})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		log.Errorf("Unable to hash the password: %v", err)
		return existedUser,
			echo.NewHTTPError(http.StatusInternalServerError, model.Response{Msg: "Unable to process the password"})
	}
	user.Password = string(hashedPassword)
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		log.Errorf("Unable to insert the user :%+v", err)
		return existedUser,
			echo.NewHTTPError(http.StatusInternalServerError, model.Response{Msg: "Unable to create the user"})
	}
	return user, nil
}

func AuthenticateUser(ctx context.Context, collection dbInterface.CollectionAPI, reqUser model.User) (model.User, *echo.HTTPError) {
	var storedUser model.User //user in db
	// check whether the user exists or not
	res := collection.FindOne(ctx, bson.M{"username": reqUser.Username})
	err := res.Decode(&storedUser)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Errorf("Unable to decode retrieved user: %v", err)
		return storedUser,
			echo.NewHTTPError(http.StatusUnprocessableEntity, model.Response{Msg: "Unable to decode retrieved user"})
	}
	if err == mongo.ErrNoDocuments {
		log.Errorf("User %s does not exist.", reqUser.Username)
		return storedUser,
			echo.NewHTTPError(http.StatusBadRequest, model.Response{Msg: "User does not exist"})
	}
	//validate the password
	if !isCredValid(reqUser.Password, storedUser.Password) {
		return storedUser,
			echo.NewHTTPError(http.StatusUnauthorized, model.Response{Msg: "Credentials invalid"})
	}
	return storedUser, nil
}

func isCredValid(givenPwd, storedPwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(storedPwd), []byte(givenPwd)); err != nil {
		return false
	}
	return true
}
