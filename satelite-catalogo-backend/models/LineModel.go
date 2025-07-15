package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Line struct {
	NumeroLinea int                `bson:"numero_linea" json:"numero_linea"`
	ProductID   primitive.ObjectID `bson:"product_id" json:"product_id"`
	Cantidad    int                `bson:"cantidad" json:"cantidad"`
	Importe     float64            `bson:"importe_linea" json:"importe_linea"`
	UM          string             `bson:"um" json:"um"`
	Comentario  string             `bson:"comentario" json:"comentario"`
}
