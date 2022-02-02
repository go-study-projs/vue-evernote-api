package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title      string             `json:"title" bson:"title" validate:"required,min=1,max=30"`
	Content    string             `json:"content" bson:"content" validate:"required,min=1,max=8000"`
	UserId     primitive.ObjectID `json:"userId" bson:"user_id"`
	NotebookId primitive.ObjectID `json:"notebookId" bson:"notebook_id"`
	IsDeleted  bool               `json:"isDeleted" bson:"is_deleted"`
	CreatedAt  time.Time          `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  time.Time          `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}
