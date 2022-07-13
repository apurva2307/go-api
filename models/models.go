package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Employee struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name" validate:"required,min=3,max=32"`
	User primitive.ObjectID `json:"user,omitempty"`
}

type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name"`
	Email     string             `json:"email" validate:"required,email"`
	Password  []byte             `json:"-" validate:"required"`
	CreatedAt time.Time          `json:"createdAt"`
	OtherInfo map[string]any     `json:"otherinfo,omitempty" bson:"otherinfo,omitempty"`
}
