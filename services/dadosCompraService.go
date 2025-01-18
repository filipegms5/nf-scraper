package services

import (
	"crypto/tls"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/filipegms5/nf-scraper/models"
	"golang.org/x/net/html"
)

var dadosCompra models.DadosCompra

func FetchDadosCompra(url string) (models.DadosCompra, error) {
	doc, err := fetch(url)
	if err != nil {
		return models.DadosCompra{}, err
	}
	dadosCompra = models.DadosCompra{}
	scrapeAll(doc)
	for _, produto := range dadosCompra.Produtos {
		quantidade, err := strconv.Atoi(produto.Quantidade)
		if err == nil {
			dadosCompra.QuantidadeTotal += quantidade
		} else {
			dadosCompra.QuantidadeTotal += 1
		}
	}
	return dadosCompra, err
}

// Fetches and parses the HTML document
func fetch(url string) (*html.Node, error) {
	// Create a custom HTTP client with TLS configuration to skip SSL verification
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: customTransport}

	// Send a GET request
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func scrapeAll(n *html.Node) {

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {

		storeName(n)

		// Store address
		storeAdress(n)

		// Date
		date(n)

		// Products
		products(n)

		// Sale info
		saleInfo(n)

		// Traverse children
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
}

func storeName(n *html.Node) {
	// Store name
	if n.Type == html.ElementNode && n.Data == "h4" {
		for b := n.FirstChild; b != nil; b = b.NextSibling {
			if b.Type == html.ElementNode && b.Data == "b" && b.FirstChild != nil {
				dadosCompra.Loja = b.FirstChild.Data
			}
		}
	}
}

func date(n *html.Node) {
	dateRegex := regexp.MustCompile(`\d{2}/\d{2}/\d{4} \d{2}:\d{2}:\d{2}`)
	if n.Type == html.TextNode {
		if data := strings.TrimSpace(n.Data); dateRegex.MatchString(data) {
			dadosCompra.Data = data
		}
	}
}

func storeAdress(n *html.Node) {
	// Store address
	if n.Type == html.ElementNode && n.Data == "td" {
		for _, attr := range n.Attr {
			if attr.Key == "style" && attr.Val == "border-top: 0px; display: block; font-style: italic;" && n.FirstChild != nil {
				dadosCompra.Endereco = n.FirstChild.Data
			}
		}
	}
}

func products(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "tr" {
		produto := models.Produto{}
		count := 0
		for td := n.FirstChild; td != nil; td = td.NextSibling {
			if td.Type == html.ElementNode && td.Data == "td" {
				text := ""
				if td.FirstChild != nil {
					if td.FirstChild.Type == html.ElementNode && td.FirstChild.Data == "h7" {
						text = td.FirstChild.FirstChild.Data
						produto.Nome = strings.Split(text, "\n")[0]
					} else {
						text = td.FirstChild.Data
						switch count {
						case 1:
							parts := strings.Split(text, ":")
							if len(parts) > 1 {
								produto.Quantidade = strings.TrimSpace(parts[1])
							}
						case 2:
							parts := strings.Split(text, ":")
							if len(parts) > 1 {
								produto.Unidade = strings.TrimSpace(parts[1])
							}
						case 3:
							parts := strings.Split(text, ": R$ ")
							if len(parts) > 1 {
								produto.Valor = strings.TrimSpace(parts[1])
							}
						}
					}
				}
				count++
			}
		}
		if produto.Nome != "" {
			dadosCompra.Produtos = append(dadosCompra.Produtos, produto)
		}
	}
}

func saleInfo(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "strong" && n.FirstChild != nil {
		text := n.FirstChild.Data
		if strings.Contains(text, ".") {
			dadosCompra.ValorTotal = text
		} else if n.FirstChild.Type == html.ElementNode && n.FirstChild.Data == "div" {
			dadosCompra.FormaPagamento = n.FirstChild.FirstChild.Data
		}
	}
}
