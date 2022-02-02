package handler

import (
	"github.com/go-playground/validator/v10"
)

var (
	v = validator.New()
)

type userValidator struct {
	validator *validator.Validate
}

func (u *userValidator) Validate(i interface{}) error {
	return u.validator.Struct(i)
}

type notebookValidator struct {
	validator *validator.Validate
}

func (nb *notebookValidator) Validate(i interface{}) error {
	return nb.validator.Struct(i)
}

type noteValidator struct {
	validator *validator.Validate
}

func (n *noteValidator) Validate(i interface{}) error {
	return n.validator.Struct(i)
}
