package models

type DadosCompra struct {
	Produtos       Produtos
	Loja           string
	Endereco       string
	FormaPagamento string
	Data           string
	ValorTotal     string
}
