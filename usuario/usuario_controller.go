package usuario

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UsuarioController struct {
	usuarioService *UsuarioService
}

func NewUsuarioController(usuarioService *UsuarioService) *UsuarioController {
	return &UsuarioController{
		usuarioService: usuarioService,
	}
}

func (me *UsuarioController) Save(c *gin.Context) {

	if c.GetHeader("X-NA-C") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cliente é obrigatório"})
	}

	var usuario Usuario

	if err := c.Bind(&usuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Não foi possível formatar o usuario"})
		return
	}

	usuario.ClientID = c.GetHeader("X-NA-C")
	//TODO - LINK COM TOKEN PARA PRIMEIRO ACESSO ATIVAR
	usuario.Status = "ACTIVE"

	if err := me.usuarioService.Salvar(&usuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Não foi possível salvar o usuario"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"result": true})
}

func (me *UsuarioController) ListUsers(c *gin.Context) {

	if c.GetHeader("X-NA-C") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cliente é obrigatório"})
	}

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
	usuarios, count, err := me.usuarioService.GetUsersByClient(c.GetHeader("X-NA-C"), from, size)

	if err != nil {
		msg := fmt.Sprintf("Erro ao recuperar usuários,  %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": msg})
	}

	result := UsuarioResult{
		Total:  int64(count),
		Result: usuarios,
	}

	for i := 0; i < len(result.Result); i++ {
		result.Result[i].Password = ""
	}

	c.JSON(http.StatusOK, result)
}

func (me *UsuarioController) GetUserById(c *gin.Context) {

	if c.Param("id") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Usuário é obrigatório"})
	}

	result, err := me.usuarioService.GetUserById(c.Query("user"))

	if err != nil {
		msg := fmt.Sprintf("Erro ao recuperar usuario,  %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": msg})
	}

	c.JSON(http.StatusOK, result)
}

func (me *UsuarioController) Login(c *gin.Context) {

	var usuario Usuario

	if err := c.Bind(&usuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Não foi possível formatar o usuario"})
	}

	user, err := me.usuarioService.Login(&usuario)
	if err != nil {
		msg := fmt.Sprintf("Erro ao realizar login,  %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}

	user.Password = ""

	c.Header("X-Request-Id", "123")

	c.JSON(http.StatusOK, gin.H{"result": user})
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
