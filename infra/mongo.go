package infra

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/mgo.v2"
)

var session *mgo.Session

func NewSession() (*mgo.Session, error) {
	mongoUrl := os.Getenv("MONGO_URL")

	if mongoUrl == "" {
		return nil, errors.New("Mongo Url is not defined")
	}

	var err error
	if session, err = mgo.Dial(mongoUrl); err != nil {
		return nil, errors.New(fmt.Sprintf("Error to connect mongoDd Url is not defined: %s", err))
	}

	session.SetMode(mgo.Monotonic, true)

	return session, nil
}

func RefreshSession() {
	if session == nil {
		NewSession()
	} else {
		session.Refresh()
	}

	return
}

func NewUsuarioCollection(database string) (*mgo.Collection, error) {
	if session == nil {

		if s, err := NewSession(); err == nil {
			session = s
		} else {
			return nil, err
		}
	}

	newCollection := session.DB(database).C("users")

	idx_username := mgo.Index{
		Name:       "idx_username",
		Key:        []string{"userName"},
		Background: true,
		Sparse:     true,
		Unique:     true,
	}

	if err := newCollection.EnsureIndex(idx_username); err != nil {
		return nil, fmt.Errorf("error on MongoDB ensure mapping index: %q", err)
	}

	return newCollection, nil
}

func NewClienteCollection(database string) (*mgo.Collection, error) {
	if session == nil {

		if s, err := NewSession(); err == nil {
			session = s
		} else {
			return nil, err
		}
	}

	newCollection := session.DB(database).C("clients")

	idx_cnpj := mgo.Index{
		Name:       "idx_cnpj",
		Key:        []string{"cnpj"},
		Background: true,
		Sparse:     true,
		Unique:     true,
	}

	if err := newCollection.EnsureIndex(idx_cnpj); err != nil {
		return nil, fmt.Errorf("error on MongoDB ensure mapping index: %q", err)
	}

	return newCollection, nil
}

func NewFilesCollection(database string) (*mgo.Collection, error) {
	if session == nil {

		if s, err := NewSession(); err == nil {
			session = s
		} else {
			return nil, err
		}
	}

	newCollection := session.DB(database).C("files")

	if err := newCollection.EnsureIndex(mgo.Index{Key: []string{"client.id"}}); err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot create index [client.id] to files collection: %s", err))
	}

	if err := newCollection.EnsureIndex(mgo.Index{Key: []string{"client.id", "user.id"}}); err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot create index [client.id, user.id] to files collection: %s", err))
	}

	if err := newCollection.EnsureIndex(mgo.Index{Key: []string{"-createdAt"}}); err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot create index [createdAt] to files collection: %s", err))
	}

	if err := newCollection.EnsureIndex(mgo.Index{Key: []string{"client.id", "user.id", "-createdAt"}}); err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot create index [client.id, user.id, createdAt] to files collection: %s", err))
	}

	return newCollection, nil
}