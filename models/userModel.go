package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// User is the model that governs all notes object retrieved or inserted into the DB
type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	Company       *string            `json:"company" validate:"required"`
	Nit_company   *string            `json:"nit_company" validate:"required"`
	Password      *string            `json:"password" validate:"required,min=6"`
	Email         *string            `json:"email" validate:"email,required"`
	Phone         *string            `json:"phone" validate:"required"`
	Token         *string            `json:"token"`
	Refresh_token *string            `json:"refresh_token"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"user_id"`
	Files         []Files            `json:"files"`
}
