package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CC struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Numero int                `bson:"numero,omitempty" json:"numero"`
	Nombre string             `bson:"nombre" json:"nombre"`
	Jefe   primitive.ObjectID `bson:"jefe,omitempty" json:"jefe,omitempty"`
}
