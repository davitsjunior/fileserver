package upfile

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type (
	File struct {
		ID           bson.ObjectId `json:"id,omitempty" bson:"_id"`
		ClientId     string        `json:"clientId,omitempty" bson:"clientId"`
		UserId       string        `json:"userId,omitempty" bson:"userId"`
		FileName     string        `json:"fileName" bson:"fileName"`
		CreatedAt    time.Time     `json:"createdAt" bson:"createdAt"`
		InternalName string        `json:"internalName,omitempty" bson:"internalName,omitempty"`
		DocYear		 string			`json:"docYear" bson:"docYear"`
		Edicao		 int			`json:"edicao" bson:"edicao"`
		DocText		 string			`json:"docText" bson:"docText"`
	}

	FileResult struct {
		Files []File `json:"files"`
		Total int    `json:"total"`
	}
)
