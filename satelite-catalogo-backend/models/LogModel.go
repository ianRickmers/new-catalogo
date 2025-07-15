package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RequestLog struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	RequestID     primitive.ObjectID `json:"request_id" bson:"request_id"`                             // ID del RequestModel asociado
	Timestamp     time.Time          `json:"timestamp" bson:"timestamp"`                               // Fecha y hora del evento
	EventType     string             `json:"event_type" bson:"event_type"`                             // Tipo de evento (creación, actualización, cambio de estado, validación, error, etc.)
	Description   string             `json:"description" bson:"description"`                           // Descripción detallada del evento
	PreviousState *Solicitud         `json:"previous_state,omitempty" bson:"previous_state,omitempty"` // (opcional) Estado previo de la solicitud (para cambios importantes)
	NewState      *Solicitud         `json:"new_state,omitempty" bson:"new_state,omitempty"`           // (opcional) Nuevo estado de la solicitud (para cambios importantes)
	UserID        primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`                         // (opcional) Usuario que realizó el cambio (analista, solicitante, sistema, etc.)
}
