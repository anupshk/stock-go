package cli

import "go.mongodb.org/mongo-driver/bson/primitive"

type Share struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Client       primitive.ObjectID `bson:"client,omitempty" json:"client,omitempty"`
	Symbol       string             `bson:"symbol" json:"symbol"`
	Balance      float32            `bson:"balance" json:"balance"`
	PCP          float32            `bson:"prev_close_price" json:"prev_close_price"`
	LTP          float32            `bson:"last_tran_price" json:"last_tran_price"`
	ImportedDate string             `bson:"imported_date" json:"imported_date"`
	ImportedAt   primitive.DateTime `bson:"imported_at" json:"imported_at"`
}
