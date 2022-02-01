package handler

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-study-projs/vue-evernote-api/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-study-projs/vue-evernote-api/dao"
	"github.com/go-study-projs/vue-evernote-api/dbInterface"
	"github.com/go-study-projs/vue-evernote-api/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type UsersHandler struct {
	Col dbInterface.CollectionAPI
}

var (
	usernameRegex = regexp.MustCompile("^[\\w\u4e00-\u9fa5]{3,15}$")
)

func (h *UsersHandler) CreateUser(c echo.Context) error {
	var user model.User
	c.Echo().Validator = &userValidator{validator: v}
	if err := c.Bind(&user); err != nil {
		log.Errorf("Unable to bind to user struct.")
		return utils.Json(c, http.StatusBadRequest, "Unable to parse the request payload.")
	}

	if err := c.Validate(user); err != nil {
		log.Errorf("Unable to validate the requested body.")
		return utils.Json(c, http.StatusBadRequest, "Unable to validate request body.")
	}

	if !usernameRegex.MatchString(user.Username) {
		log.Errorf("userName %s is not valid.", user.Username)
		return utils.Json(c, http.StatusBadRequest, fmt.Sprintf("userName %s is not valid.", user.Username))
	}

	resUser, httpError := dao.InsertUser(context.Background(), user, h.Col)
	if httpError != nil {
		return httpError
	}

	return utils.Json(c, http.StatusCreated, "注册成功", model.User{
		ID:        resUser.ID,
		Username:  resUser.Username,
		CreatedAt: resUser.CreatedAt,
		UpdatedAt: resUser.UpdatedAt,
	})

}

func (h *UsersHandler) AuthnUser(c echo.Context) error {
	var user model.User
	c.Echo().Validator = &userValidator{validator: v}
	if err := c.Bind(&user); err != nil {
		log.Errorf("Unable to bind to user struct.")
		return c.JSON(http.StatusUnprocessableEntity,
			model.ErrorMessage{Message: "Unable to parse the request payload."})
	}
	if err := c.Validate(user); err != nil {
		log.Errorf("Unable to validate the requested body.")
		return c.JSON(http.StatusBadRequest,
			model.ErrorMessage{Message: "Unable to validate request payload"})
	}
	user, httpError := dao.AuthenticateUser(context.Background(), user, h.Col)
	if httpError != nil {
		return c.JSON(httpError.Code, httpError.Message)
	}
	token, err := user.CreateToken()
	if err != nil {
		log.Errorf("Unable to generate the token.")
		return c.JSON(http.StatusInternalServerError,
			model.ErrorMessage{Message: "Unable to generate the token"})
	}
	c.Response().Header().Set("x-auth-token", "Bearer "+token)
	return c.JSON(http.StatusOK, model.User{Username: user.Username})
}

func (h *UsersHandler) Login(c echo.Context) error {
	var user model.User
	c.Echo().Validator = &userValidator{validator: v}
	if err := c.Bind(&user); err != nil {
		log.Errorf("Unable to bind to user struct.")
		return utils.Json(c, http.StatusBadRequest, "Unable to parse the request payload.")
	}

	if err := c.Validate(user); err != nil {
		log.Errorf("Unable to validate the requested body.")
		return utils.Json(c, http.StatusBadRequest, "Unable to parse the request payload.")
	}

	if !usernameRegex.MatchString(user.Username) {
		log.Errorf("userName %s is not valid.", user.Username)
		return utils.Json(c, http.StatusBadRequest, fmt.Sprintf("userName %s is not valid.", user.Username))
	}

	var existedUser model.User
	res := dao.FindOneUser(context.Background(),
		h.Col,
		bson.M{"username": user.Username, "password": user.Password})
	err := res.Decode(&existedUser)

	if err != nil && err != mongo.ErrNoDocuments {
		log.Errorf("Unable to decode retrieved user: %v", err)
		return utils.Json(c, http.StatusUnprocessableEntity, "Unable to decode retrieved user")
	}

	token, err := user.CreateToken()
	if err != nil {
		log.Errorf("Unable to generate the token.")
		return utils.Json(c, http.StatusInternalServerError, "Unable to generate the token")

	}
	//c.Response().Header().Set("x-auth-token", "Bearer "+token)
	return utils.Json(c, http.StatusOK, "登录成功", map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"token":    token,
	})
}
