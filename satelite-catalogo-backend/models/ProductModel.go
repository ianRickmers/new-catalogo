package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Number             string             `bson:"number,omitempty" json:"number"`
	Descripcion        string             `bson:"descripcion" json:"descripcion"`
	Licitacion         string             `bson:"licitacion,omitempty" json:"licitacion,omitempty"`
	IDConvenio         string             `bson:"id_convenio,omitempty" json:"id_convenio,omitempty"`
	NombreProveedor    string             `bson:"nombre_proveedor,omitempty" json:"nombre_proveedor,omitempty"`
	RutProveedor       string             `bson:"rut_proveedor,omitempty" json:"rut_proveedor,omitempty"`
	IDProduct          string             `bson:"id_product,omitempty" json:"id_product,omitempty"`
	Region             string             `bson:"region,omitempty" json:"region,omitempty"`
	Marca              string             `bson:"marca,omitempty" json:"marca,omitempty"`
	Modelo             string             `bson:"modelo,omitempty" json:"modelo,omitempty"`
	Precio             float64            `bson:"precio,omitempty" json:"precio,omitempty"`
	FechaActualizacion time.Time          `bson:"fecha_actualizacion,omitempty" json:"fecha_actualizacion,omitempty"`
	ConvenioMarco      string             `bson:"convenio_marco,omitempty" json:"convenio_marco,omitempty"`
	UM                 string             `bson:"UM,omitempty" json:"UM,omitempty"`
	Categoria          string             `bson:"categoria,omitempty" json:"categoria,omitempty"`
	IDCategoria        string             `bson:"id_categoria,omitempty" json:"id_categoria,omitempty"`
}
