package main

import (
	"log"
	"os"

	"time"

	"fileserver/upfile"
	"fileserver/usuario"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	checkEnvs()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-CLIENT", "X-USER"},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}))

	/*USUARIO*/
	usuarioDAO, err := usuario.NewUsuarioDAO()
	if err != nil {
		log.Printf("Erro ao instanciar Usuario DAO ERR: %s", err.Error())
	}
	usuarioService := usuario.NewUsuarioService(usuarioDAO)
	usuarioController := usuario.NewUsuarioController(usuarioService)
	router.POST("/user", usuarioController.Save)
	router.POST("/user/auth", usuarioController.Login)
	router.GET("/user/:id", usuarioController.GetUserById)
	router.GET("/user", usuarioController.ListUsers)
	/*USUARIO*/

	/*ARQUIVO*/
	fileService, _ := upfile.NewFilesService()
	uploaderController := upfile.NewUploaderController(fileService)
	router.POST("/upload", uploaderController.Upload)
	router.GET("/files/list", uploaderController.GetFiles)
	router.GET("/allfiles/list", uploaderController.GetAllFiles)
	/*ARQUIVO*/

	const DOWNLOADS_PATH = ""

	router.GET("/download/:filename", func(ctx *gin.Context) {
		fileName := ctx.Param("filename")
		targetPath := filepath.Join(DOWNLOADS_PATH, fileName)
		//This ckeck is for example, I not sure is it can prevent all possible filename attacks - will be much better if real filename will not come from user side. I not even tryed this code
		if !strings.HasPrefix(filepath.Clean(targetPath), DOWNLOADS_PATH) {
			ctx.String(403, "Look like you attacking me")
			return
		}
		fmt.Println(targetPath)
		//Seems this headers needed for some browsers (for example without this headers Chrome will download files as txt)
		ctx.Header("Content-Description", "File Transfer")
		ctx.Header("Content-Transfer-Encoding", "binary")
		ctx.Header("Content-Disposition", "attachment; filename="+fileName)
		ctx.Header("Content-Type", "application/octet-stream")
		ctx.File(targetPath)
	})

	port := os.Getenv("SERVER_PORT")
	router.Run(":" + port)
}

func checkEnvs() {
	checkRequired("SERVER_PORT", "MONGO_URL", "DIR_LOCAL")
}

func checkRequired(envVarArgs ...string) {
	for _, envVar := range envVarArgs {
		if os.Getenv(envVar) == "" {
			log.Fatal("Environment variable '%s' is required.", envVar)
			continue
		}
		log.Printf("Environment variable '%s' is ok with value: %s", envVar, os.Getenv(envVar))
	}
}
