package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CC = "Centro de Costo"

type User struct {
	ID        primitive.ObjectID   `json:"_id"   bson:"_id,omitempty"`
	Username  string               `json:"username,omitempty"   bson:"username,omitempty"`
	Email     string               `json:"email,omitempty"   bson:"email,omitempty"`
	Rut       string               `json:"rut,omitempty"   bson:"rut,omitempty"`
	Role      []Role               `json:"role,omitempty"   bson:"role,omitempty"`
	CreatedAt primitive.DateTime   `json:"created_at,omitempty" bson:"created_at,omitempty" swaggertype:"string"`
	CC        []primitive.ObjectID `bson:"cc" json:"cc"`
}

type Role string

const (
	USER  Role = "Usuario"
	ADMIN Role = "Administrador"
)

func NewUser(username string, email string, rut string, role []Role) *User {
	return &User{
		Username:  username,
		Email:     email,
		Role:      role,
		Rut:       rut,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		//CC:        cc,
	}
}
