package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Solicitud struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CC              primitive.ObjectID `bson:"cc" json:"cc"`
	Lines           []Line             `bson:"lines" json:"lines"`
	Solicitante     primitive.ObjectID `bson:"solicitante" json:"solicitante"`
	Aprobador       primitive.ObjectID `bson:"aprobador" json:"aprobador"`
	Description     string             `bson:"description" json:"description"`
	Documents       []string           `bson:"documents" json:"documents"`
	State           string             `bson:"state" json:"state"`
	FechaSolicitud  time.Time          `bson:"fecha_solicitud" json:"fecha_solicitud"`
	FechaContable   time.Time          `bson:"fecha_contable" json:"fecha_contable"`
	Moneda          string             `bson:"moneda" json:"moneda"`
	NombreSolicitud string             `bson:"nombre_solicitud" json:"nombre_solicitud"`
	ImporteTotal    float64            `bson:"importe_total" json:"importe_total"`
}
