package handler

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

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
		return c.JSON(http.StatusBadRequest,
			model.ErrorMessage{Message: "Unable to parse the request payload."})
	}

	if err := c.Validate(user); err != nil {
		log.Errorf("Unable to validate the requested body.")
		return c.JSON(http.StatusBadRequest,
			model.ErrorMessage{Message: "Unable to validate request body"})
	}

	if !usernameRegex.MatchString(user.Username) {
		log.Errorf("userName %s is not valid.", user.Username)
		return c.JSON(http.StatusBadRequest,
			model.ErrorMessage{Message: fmt.Sprintf("userName %s is not valid.", user.Username)})
	}

	resUser, httpError := dao.InsertUser(context.Background(), user, h.Col)
	if httpError != nil {
		return c.JSON(httpError.Code, httpError.Message)
	}
	token, err := user.CreateToken()
	if err != nil {
		log.Errorf("Unable to generate the token.")
		return echo.NewHTTPError(http.StatusInternalServerError,
			model.ErrorMessage{Message: "Unable to generate the token"})
	}
	c.Response().Header().Set("x-auth-token", "Bearer "+token)
	return c.JSON(http.StatusCreated, resUser)
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
