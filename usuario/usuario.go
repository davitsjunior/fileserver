package usuario

import (
	"math/rand"
	"fileserver/permissao"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type (
	Usuario struct {
		ID         bson.ObjectId         `json:"id" bson:"_id"`
		ClientID   string                `json:"clientId" bson:"clientId"`
		Nome       string                `json:"nome" bson:"nome"`
		UserName   string                `json:"userName" bson:"userName"`
		Password   string                `json:"password" bson:"password"`
		Token      string                `json:"token" bson:"token"`
		Permissoes []permissao.Permissao `json:"permissoes" bson:"permissoes"`
		CreatedAt  time.Time             `json:"createdAt" bson:"createdAt"`
		UpdatedAt  time.Time             `json:"updatedAt" bson:"updatedAt"`
		Status     string                `json:"status" bson:"status"`
		Ramal      string                `json:"ramal" bson:"ramal"`
		FixPermissao string				`json:"fixPermissao" bson:"fixPermissao"`
	}

	UsuarioResult struct {
		Total  int64     `json:"total"`
		Result []Usuario `json:"result"`
	}
)

func (me *Usuario) GerarSenhaRandomica() {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	me.Password = string(b)
}

func (me *Usuario) ValidarUsuario(usuario *Usuario) bool {

	if usuario.Status == "ACTIVE" && me.Password == usuario.Password {
		return true
	}
	return false
}
