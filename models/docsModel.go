package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is the model that governs all notes object retrieved or inserted into the DB
type Doc struct {
	ID          primitive.ObjectID `bson:"_id"`
	User_name   *string            `json:"user_name" validate:"required,min=2,max=100"`
	Title       *string            `json:"title" validate:"required,min=2,max=100"`
	Description *string            `json:"description" validate:"required,min=6"`
}
