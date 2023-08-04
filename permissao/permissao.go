package permissao

import "gopkg.in/mgo.v2/bson"

type (
	Permissao struct {
		ID    bson.ObjectId `json:"id" bson:"_id"`
		Nome  string        `json:"nome" bson:"_id"`
		Menus []Menu        `json:"menus" bson:"menus"`
	}

	Menu struct {
		ID              string           `json:"id" bson:"_id"`
		Nome            string           `json:"nome" bson:"_id"`
		Funcionalidades []Funcionalidade `json:"funcionalidades" bson:"funcionalidades"`
	}

	Funcionalidade struct {
		ID   string `json:"id" bson:"_id"`
		Nome string `json:"nome" bson:"_id"`
	}
)
