package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id"  form:"id,omitempty"`
	Username string             `json:"username" form:"username"`
	Password string             `json:"password" form:"password" binding:"required"`
}
