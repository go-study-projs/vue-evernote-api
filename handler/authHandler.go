package handler

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-study-projs/vue-evernote-api/utils"

	"github.com/go-study-projs/vue-evernote-api/dao"
	"github.com/go-study-projs/vue-evernote-api/dbInterface"
	"github.com/go-study-projs/vue-evernote-api/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type UserHandler struct {
	Col dbInterface.CollectionAPI
}

var (
	// 用户名，长度3~15个字符，仅限于字母数字下划线中文
	usernameRegex = regexp.MustCompile("^[\\w\u4e00-\u9fa5]{3,15}$")
)

func (h *UserHandler) CreateUser(c echo.Context) error {
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

	resUser, httpError := dao.InsertUser(context.Background(), h.Col, user)
	if httpError != nil {
		return httpError
	}

	token, err := utils.CreateToken(resUser)
	if err != nil {
		log.Errorf("Unable to generate the token.")
		return utils.Json(c, http.StatusInternalServerError, "Unable to generate the token")

	}
	return utils.Json(c, http.StatusOK, "注册成功", map[string]interface{}{
		"id":        resUser.ID,
		"username":  resUser.Username,
		"token":     token,
		"createdAt": resUser.CreatedAt,
		"updatedAt": resUser.UpdatedAt,
	})
}

func (h *UserHandler) Login(c echo.Context) error {
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

	storedUser, httpError := dao.AuthenticateUser(context.Background(), h.Col, user)
	if httpError != nil {
		return httpError
	}

	token, err := utils.CreateToken(storedUser)
	if err != nil {
		log.Errorf("Unable to generate the token.")
		return utils.Json(c, http.StatusInternalServerError, "Unable to generate the token")

	}
	return utils.Json(c, http.StatusOK, "登录成功", map[string]interface{}{
		"id":        storedUser.ID,
		"username":  storedUser.Username,
		"token":     token,
		"createdAt": storedUser.CreatedAt,
		"updatedAt": storedUser.UpdatedAt,
	})
}
