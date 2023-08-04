package cliente

type (
	Cliente struct {
		ID         string       		`json:"id" bson:"_id"`
		CNPJ	   string 				`json:"cnpj" bson:"cnpj"`
		Nome       string                `json:"nome" bson:"nome"`
		Ativo		bool			`json:"ativo" bson:"ativo"`
		ZapAccounts []ZapAccount	`json:"zapAccounts" bson:"zapAccounts"`
	}

	ClienteResult struct {
		Total  int64     `json:"total"`
		Result []Cliente `json:"result"`
	}

	ZapAccount struct {
		Numero string `json:"numero" bson:"numero"`
		Key		string `json:"key" bson:"key"`
		Mailing bool	`json:"mailing"  bson:"mailing"`
		Atendimento bool `json:"atendimento" bson:"atendimento"`

	}
)
