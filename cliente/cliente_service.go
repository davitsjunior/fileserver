package cliente

import (
	"errors"
)

type ClienteService struct {
	clienteDAO   *ClienteDAO
}

func NewClienteService(clienteDAO *ClienteDAO) *ClienteService {
	return &ClienteService{
		clienteDAO:   clienteDAO,
	}
}

func (me *ClienteService) Salvar(cliente *Cliente) error {

	err := me.clienteDAO.Save(cliente)

	if err != nil {
		return err
	}

	return nil
}

func (me *ClienteService) GetById(clientId string) (*Cliente, error) {

	return me.clienteDAO.GetClienteById(clientId)
}

func (me *ClienteService) GetZapKeyByNumber(clientId, numero string) (ZapAccount, error) {

	cliente, err := me.clienteDAO.GetClienteById(clientId)
	if err != nil {
		return ZapAccount{}, err
	}

	if len(cliente.ZapAccounts) < 1 {
		return ZapAccount{}, errors.New("Cliente não possui zap account cadastrado")
	}

	for _, v := range cliente.ZapAccounts {
		if v.Numero == numero {
			return v, nil
		}
	}

	return ZapAccount{}, errors.New("Não localizamos o numero no cadastro zap do cliente")
}

func(me *ClienteService) LoadClientsByStatus(ativo bool) ([]Cliente, error) {

	return me.clienteDAO.LoadClientsByStatus(ativo)
}