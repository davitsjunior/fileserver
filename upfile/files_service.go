package upfile

import (
	"errors"
	"fmt"
	"fileserver/infra"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type (
	FilesService struct {
		Collection *mgo.Collection
	}
)

func NewFilesService() (*FilesService, error) {
	return NewFilesServiceWithDBName("fileserver")
}

func NewFilesServiceWithDBName(dbName string) (*FilesService, error) {
	if coll, err := infra.NewFilesCollection(dbName); err == nil {
		return &FilesService{
			Collection: coll,
		}, nil
	} else {
		return nil, err
	}
}

func (self *FilesService) Save(file *File) error {
	if file == nil {
		return errors.New("File mustn't be nil")
	}

	infra.RefreshSession()

	err := self.Collection.Insert(file)
	if err != nil {
		strError := fmt.Sprintf("Mongodb insert file failed: %s", err.Error())
		log.Println(strError)
		return errors.New(strError)
	}

	return nil
}

func (self *FilesService) GetFiles(from, size int) (files []File, count int, err error) {

	query := bson.M{}

	infra.RefreshSession()

	if count, err = self.Collection.Find(query).Count(); err != nil {
		return nil, 0, err
	}

	find := self.Collection.Find(&query)
	find.Skip(from)
	find.Limit(size)
	find.Sort("-createdAt")

	err = find.All(&files)
	if err != nil {
		return nil, 0, err
	}

	return files, count, nil
}

func (self *FilesService) GetById(fileId string) (file File, err error) {

	if !bson.IsObjectIdHex(fileId) {
		err = errors.New("The value" + fileId + "is a invalid ObjectId")
		return
	}

	query := bson.M{}
	query["_id"] = bson.ObjectIdHex(fileId)

	infra.RefreshSession()

	if err = self.Collection.Find(query).One(&file); err != nil {
		return
	}

	return
}

func (self *FilesService) Update(file File) (updated bool, err error) {

	infra.RefreshSession()

	if err = self.Collection.UpdateId(file.ID, file); err != nil {
		return false, err
	}

	return true, err
}

func (self *FilesService) CountByDocYear(docYear string) (int, error) {

	infra.RefreshSession()

	total, err := self.Collection.Find(bson.M{"docYear": docYear}).Count()
	if err != nil {
		return 0, err
	}

	return total, err
}