package upfile

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

type UploaderController struct {
	FileService *FilesService
}

func NewUploaderController(fileservice *FilesService) *UploaderController {
	return &UploaderController{
		FileService: fileservice,
	}
}

func (self *UploaderController) Upload(c *gin.Context) {

	var err error

	file, header, err := c.Request.FormFile("file")
	docText := c.Request.FormValue("texto")

	if err != nil || header == nil {
		fmt.Println("01")
		fmt.Println(header)
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "File load error"})
		return
	}

	//dir := os.Getenv("/data/files/teste/")
	//
	//if _, err := os.Stat(dir); err != nil {
	//	fmt.Println(err.Error())
	//}
	//if os.IsNotExist(err) {
	//	os.Mkdir(dir, 0777)
	//}

	uuid := bson.NewObjectId()
	filePath := uuid.Hex() + ".pdf"

	out, err := os.Create(filePath)
	if err != nil {
		fmt.Println("01")
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Create dir error", "error": err.Error()})
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println("02")
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Save file error", "error": err.Error()})
		return
	}

	clientId := c.GetHeader("X-CLIENT")
	userId := c.GetHeader("X-USER")

	increment, err := self.FileService.CountByDocYear("V")
	if err != nil {
		fmt.Println("03")
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	edicao := increment + 1
	history := File{
		ID:           bson.NewObjectId(),
		ClientId:     clientId,
		UserId:       userId,
		FileName:     header.Filename,
		CreatedAt:    time.Now(),
		InternalName: uuid.Hex(),
		DocYear:      "V",
		Edicao:       edicao,
		DocText:      docText,
	}

	err = self.FileService.Save(&history)
	if err != nil {
		fmt.Println("04")
		fmt.Println(err.Error())
		log.Printf("Upload file error %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": true})
}

func (self *UploaderController) GetFiles(c *gin.Context) {

	var from, size int
	var err error
	if from, err = strconv.Atoi(c.DefaultQuery("from", "0")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Query param 'from' precisa ser um numero"})
		return
	}
	if size, err = strconv.Atoi(c.DefaultQuery("size", "10")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Query param 'size' precisa ser um numero"})
		return
	}

	files, count, err := self.FileService.GetFiles(from, size)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	fileResult := FileResult{Files: files, Total: count}

	c.JSON(http.StatusOK, fileResult)
}

func (self *UploaderController) GetAllFiles(c *gin.Context) {
	files, count, err := self.FileService.GetAllFiles()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	fileResult := FileResult{Files: files, Total: count}

	c.JSON(http.StatusOK, fileResult)
}
