package models

type DadosCompra struct {
	Produtos        []Produto
	Loja            string
	Endereco        string
	FormaPagamento  string
	Data            string
	ValorTotal      string
	QuantidadeTotal int `default:"0"`
}
