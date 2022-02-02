package handler

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-study-projs/vue-evernote-api/dao"
	"github.com/go-study-projs/vue-evernote-api/dbInterface"
	"github.com/go-study-projs/vue-evernote-api/model"
	"github.com/go-study-projs/vue-evernote-api/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type NotebookHandler struct {
	Col  dbInterface.CollectionAPI
	NCol dbInterface.CollectionAPI
}

func (h *NotebookHandler) GetNotebooks(c echo.Context) error {
	userId := utils.ParseToken(c.Request().Header.Get("x-auth-token"))

	notebooks, httpError := dao.FindNotebooks(context.Background(), h.Col, userId)
	if httpError != nil {
		return c.JSON(httpError.Code, httpError.Message)
	}

	return c.JSON(http.StatusOK, notebooks)
}

func (h *NotebookHandler) CreateNotebook(c echo.Context) error {
	var notebook model.Notebook

	c.Echo().Validator = &notebookValidator{validator: v}
	if err := c.Bind(&notebook); err != nil {
		log.Errorf("Unable to bind to notebook struct.")
		return utils.Json(c, http.StatusBadRequest, "Unable to parse the request payload.")
	}

	if err := c.Validate(notebook); err != nil {
		log.Errorf("Unable to validate the requested body.")
		return utils.Json(c, http.StatusBadRequest, "Unable to validate request body.")
	}
	notebook.UserId = utils.ParseToken(c.Request().Header.Get("x-auth-token"))
	resNotebook, httpError := dao.InsertNotebook(context.Background(), notebook, h.Col)
	if httpError != nil {
		return httpError
	}

	return utils.Json(c, http.StatusCreated, "创建笔记本成功", resNotebook)
}

func (h *NotebookHandler) UpdateNoteBook(c echo.Context) error {
	var newNotebook model.Notebook

	c.Echo().Validator = &notebookValidator{validator: v}
	if err := c.Bind(&newNotebook); err != nil {
		log.Errorf("Unable to bind to notebook struct.")
		return utils.Json(c, http.StatusBadRequest, "Unable to parse the request payload.")
	}

	if err := c.Validate(newNotebook); err != nil {
		log.Errorf("Unable to validate the requested body.")
		return utils.Json(c, http.StatusBadRequest, "Unable to validate request body.")
	}
	notebookId, _ := primitive.ObjectIDFromHex(c.Param("notebookId"))
	httpError := dao.ModifyNoteBook(context.Background(), h.Col, notebookId, newNotebook)
	if httpError != nil {
		return httpError
	}
	return utils.Json(c, http.StatusOK, "修改成功")
}

func (h *NotebookHandler) DeleteNoteBook(c echo.Context) error {
	notebookId, _ := primitive.ObjectIDFromHex(c.Param("notebookId"))

	_, httpError := dao.DeleteNoteBook(context.Background(), h.Col, h.NCol, notebookId)
	if httpError != nil {
		return httpError
	}

	return utils.Json(c, http.StatusOK, "删除成功")
}
