package usuario

import (
	"errors"
	"fmt"
)

type UsuarioService struct {
	usuarioDAO   *UsuarioDAO
}

func NewUsuarioService(usuarioDAO *UsuarioDAO) *UsuarioService {
	return &UsuarioService{
		usuarioDAO:   usuarioDAO,
	}
}

func (me *UsuarioService) Salvar(usuario *Usuario) error {

	usuario.GerarSenhaRandomica()

	err := me.usuarioDAO.Save(usuario)

	if err != nil {
		return err
	}

	return nil
}

func (me *UsuarioService) GetUserById(userId string) (*Usuario, error) {
	return me.usuarioDAO.GetUserById(userId)
}

func (me *UsuarioService) GetUsersByClient(clientId string, from, size int) ([]Usuario, int, error) {
	return me.usuarioDAO.GetUsersByClient(clientId, from, size)
}

func (me *UsuarioService) Login(usuario *Usuario) (*Usuario, error) {

	usuarioRecuperado, err := me.usuarioDAO.GetUserByUserName(usuario.UserName)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if !usuario.ValidarUsuario(usuarioRecuperado) {
		return nil, errors.New("verifique os dados de acesso")
	}

	return usuarioRecuperado, nil
}
