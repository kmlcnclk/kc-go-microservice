package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	CreatedOn time.Time          `bson:"createdOn"`
	UpdatedOn time.Time          `bson:"updatedOn"`
}
