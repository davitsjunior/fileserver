package cliente

import (
	"errors"
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fileserver/infra"
)

type (
	ClienteDAO struct {
		Collection *mgo.Collection
	}
)

func NewClienteDAO() (*ClienteDAO, error) {
	return NewClientDAODBName("fileserver")
}

func NewClientDAODBName(dbName string) (*ClienteDAO, error) {
	if coll, err := infra.NewClienteCollection(dbName); err == nil {
		return &ClienteDAO{
			Collection: coll,
		}, nil
	} else {
		return nil, err
	}
}

func (me *ClienteDAO) Save(cliente *Cliente) error {
	if cliente == nil {
		return errors.New("Usuário nao pode ser nula")
	}

	infra.RefreshSession()

	cliente.ID = cliente.CNPJ
	cliente.Ativo = true

	err := me.Collection.Insert(cliente)
	if err != nil {
		strError := fmt.Sprintf("Mongodb insert cliente failed: %s", err.Error())
		log.Println(strError)
		return errors.New(strError)
	}

	return nil
}

func (me *ClienteDAO) GetClienteById(clienteId string) (*Cliente, error) {
	if clienteId == "" {
		return nil, errors.New("Id do cliente é obrigatorio")
	}

	cliente := &Cliente{}
	err := me.Collection.Find(bson.M{"_id": clienteId}).One(cliente)
	if err != nil {
		return nil, errors.New("erro ao recuperar cliente")
	}

	return cliente, nil
}

func (me *ClienteDAO) LoadClientsByStatus(ativo bool) ([]Cliente, error) {

	var clientes []Cliente
	err := me.Collection.Find(bson.M{"ativo": ativo}).All(&clientes)
	if err != nil {
		return nil, errors.New("erro ao recuperar lista de clientes")
	}

	return clientes, nil
}