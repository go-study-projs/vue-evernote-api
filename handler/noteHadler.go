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

type NoteHandler struct {
	Col dbInterface.CollectionAPI
}

func (h NoteHandler) CreateNote(c echo.Context) error {
	var note model.Note

	c.Echo().Validator = &noteValidator{validator: v}
	if err := c.Bind(&note); err != nil {
		log.Errorf("Unable to bind to note struct.")
		return utils.Json(c, http.StatusBadRequest, "Unable to parse the request payload.")
	}

	if err := c.Validate(note); err != nil {
		log.Errorf("Unable to validate the requested body.")
		return utils.Json(c, http.StatusBadRequest, "Unable to validate request body.")
	}

	note.UserId = utils.ParseToken(c.Request().Header.Get("x-auth-token"))
	notebookId, _ := primitive.ObjectIDFromHex(c.Param("notebookId"))
	note.NotebookId = notebookId

	resNote, httpError := dao.InsertNote(context.Background(), note, h.Col)
	if httpError != nil {
		return httpError
	}

	return utils.Json(c, http.StatusCreated, "创建笔记成功", resNote)
}

func (h NoteHandler) GetNotes(c echo.Context) error {
	notebookId, _ := primitive.ObjectIDFromHex(c.Param("notebookId"))

	notes, httpError := dao.FindNotes(context.Background(), h.Col, notebookId)
	if httpError != nil {
		return httpError
	}

	return utils.Json(c, http.StatusOK, "", notes)
}

func (h NoteHandler) MoveToTrash(c echo.Context) error {
	noteId, _ := primitive.ObjectIDFromHex(c.Param("noteId"))

	httpError := dao.SoftDeleteOrRevertNote(context.Background(), h.Col, noteId, dao.SoftDelete)
	if httpError != nil {
		return httpError
	}

	return utils.Json(c, http.StatusOK, "已放入回收站")
}

func (h NoteHandler) UpdateNote(c echo.Context) error {
	var note model.Note

	c.Echo().Validator = &noteValidator{validator: v}
	if err := c.Bind(&note); err != nil {
		log.Errorf("Unable to bind to note struct.")
		return utils.Json(c, http.StatusBadRequest, "Unable to parse the request payload.")
	}

	if err := c.Validate(note); err != nil {
		log.Errorf("Unable to validate the requested body.")
		return utils.Json(c, http.StatusBadRequest, "Unable to validate request body.")
	}
	noteId, _ := primitive.ObjectIDFromHex(c.Param("noteId"))
	httpError := dao.ModifyNote(context.Background(), h.Col, noteId, note)
	if httpError != nil {
		return httpError
	}
	return utils.Json(c, http.StatusOK, "修改成功")
}

func (h NoteHandler) DeleteNote(c echo.Context) error {
	noteId, _ := primitive.ObjectIDFromHex(c.Param("noteId"))

	_, httpError := dao.DeleteNote(context.Background(), h.Col, noteId)
	if httpError != nil {
		return httpError
	}

	return utils.Json(c, http.StatusOK, "删除成功")
}

func (h NoteHandler) RevertNote(c echo.Context) error {
	noteId, _ := primitive.ObjectIDFromHex(c.Param("noteId"))

	httpError := dao.SoftDeleteOrRevertNote(context.Background(), h.Col, noteId, dao.Revert)
	if httpError != nil {
		return httpError
	}

	return utils.Json(c, http.StatusOK, "已从回收站恢复")
}

func (h NoteHandler) GetNotesInTrash(c echo.Context) error {

	userId := utils.ParseToken(c.Request().Header.Get("x-auth-token"))

	notes, httpError := dao.FindNotesInTrash(context.Background(), h.Col, userId)
	if httpError != nil {
		return httpError
	}

	return utils.Json(c, http.StatusOK, "", notes)
}
