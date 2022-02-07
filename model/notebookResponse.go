package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotebookResponse struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title      string             `json:"title" bson:"title" validate:"required,min=1,max=30"`
	UserId     primitive.ObjectID `json:"userId" bson:"user_id"`
	CreatedAt  time.Time          `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  time.Time          `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
	NoteCounts int                `json:"noteCounts" bson:"note_counts"`
}
