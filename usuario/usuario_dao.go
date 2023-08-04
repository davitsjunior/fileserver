package usuario

import (
	"errors"
	"fmt"
	"log"
	"fileserver/infra"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	UsuarioDAO struct {
		Collection *mgo.Collection
	}
)

func NewUsuarioDAO() (*UsuarioDAO, error) {
	return NewUsuarioDAODBName("fileserver")
}

func NewUsuarioDAODBName(dbName string) (*UsuarioDAO, error) {
	if coll, err := infra.NewUsuarioCollection(dbName); err == nil {
		return &UsuarioDAO{
			Collection: coll,
		}, nil
	} else {
		return nil, err
	}
}

func (me *UsuarioDAO) Save(usuario *Usuario) error {
	if usuario == nil {
		return errors.New("Usuário nao pode ser nula")
	}

	infra.RefreshSession()

	usuario.ID = bson.NewObjectId()

	err := me.Collection.Insert(usuario)
	if err != nil {
		strError := fmt.Sprintf("Mongodb insert usuario failed: %s", err.Error())
		log.Println(strError)
		return errors.New(strError)
	}

	return nil
}

func (me *UsuarioDAO) GetUserById(userId string) (*Usuario, error) {
	if userId == "" {
		return nil, errors.New("Id do usuário é obrigatorio")
	}

	user := &Usuario{}
	err := me.Collection.Find(bson.M{"_id": bson.ObjectIdHex(userId)}).One(user)
	if err != nil {
		return nil, errors.New("erro ao recuperar usuario")
	}

	return user, nil
}

func (me *UsuarioDAO) GetUsersByClient(clientId string, from, size int) ([]Usuario, int, error) {
	if clientId == "" {
		return nil, 0, errors.New("Id do cliente é obrigatorio")
	}

	count, err := me.Collection.Find(bson.M{"clientId": clientId}).Count()
	if err != nil {
		return nil, 0, errors.New("erro ao recuperar lista de usuários")
	}

	var users = make([]Usuario, 0)
	err = me.Collection.Find(bson.M{"clientId": clientId}).Skip(from).Limit(size).All(&users)
	if err != nil {
		return nil, 0, errors.New("erro ao recuperar lista de usuários")
	}

	return users, count, nil
}

func (me *UsuarioDAO) GetUserByUserName(userName string) (*Usuario, error) {
	if userName == "" {
		return nil, errors.New("username do usuário é obrigatorio")
	}

	count, err := me.Collection.Find(bson.M{"userName": userName}).Count()
	if err != nil {
		return nil, errors.New("erro ao recuperar usuario")
	}
	if count < 1 {
		return nil, errors.New("usuario não cadatrado")
	}
	user := &Usuario{}
	err = me.Collection.Find(bson.M{"userName": userName}).One(user)
	if err != nil {
		return nil, errors.New("erro ao recuperar usuario")
	}

	return user, nil
}
