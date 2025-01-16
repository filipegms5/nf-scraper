package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/filipegms5/nf-scraper/models"
	"golang.org/x/net/html"
)

var produtos models.Produtos
var dadosCompra models.DadosCompra

// Fetches and parses the HTML document
func fetch(url string) (*html.Node, error) {
	// Send a GET request
	resp, err := http.Get(url)
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

// Extracts the titles from a webpage by traversing the HTML nodes
func scrapeTitles(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "tr" { // Adjust this for different tags or classes
		//fmt.Println("chegou aqui")

		for _, attr := range n.Attr {
			if attr.Key == "class" && attr.Val == "" { // Modify as needed
				fmt.Println(n.FirstChild.Data)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		scrapeTitles(c)
	}
}

func scrapeProducts(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "tbody" && len(n.Attr) > 0 {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "tr" {
				var produto models.Produto
				var count = 0
				for td := c.FirstChild; td != nil; td = td.NextSibling {
					if td.Type == html.ElementNode && td.Data == "td" {
						if td.FirstChild != nil && td.FirstChild.Type == html.ElementNode && td.FirstChild.Data == "h7" {
							produto.Nome = td.FirstChild.FirstChild.Data
						} else {
							if count == 1 {
								produto.Quantidade = strings.Split(td.FirstChild.Data, ":")[1]
							} else if count == 2 {
								produto.Unidade = strings.Split(td.FirstChild.Data, ":")[1]
							} else if count == 3 {
								produto.Valor = strings.Split(td.FirstChild.Data, ": R$ ")[1]
							}
						}
						count++
					}
				}
				produtos.Produtos = append(produtos.Produtos, produto)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		scrapeProducts(c)
	}
}
func scrapeSaleInfo(n *html.Node) {
	count := 0
	var printStrongTags func(*html.Node)
	printStrongTags = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "strong" {
			count++
			if count == 6 {
				dadosCompra.ValorTotal = n.FirstChild.Data
			} else if count == 8 {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.ElementNode && c.Data == "div" {
						dadosCompra.FormaPagamento = c.FirstChild.Data
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			printStrongTags(c)
		}
	}
	printStrongTags(n)
}

func scrapeStoreName(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "th" {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "h4" {
				for b := c.FirstChild; b != nil; b = b.NextSibling {
					if b.Type == html.ElementNode && b.Data == "b" {
						dadosCompra.Loja = b.FirstChild.Data
					}
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		scrapeStoreName(c)
	}
}

func scrapeStoreAddress(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "td" {
		for _, attr := range n.Attr {
			if attr.Key == "style" && attr.Val == "border-top: 0px; display: block; font-style: italic;" {
				if n.FirstChild != nil {
					dadosCompra.Endereco = n.FirstChild.Data
					return
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		scrapeStoreAddress(c)
	}
}

func scrapeDate(n *html.Node) {
	if n.Type == html.TextNode {
		data := strings.TrimSpace(n.Data)
		matched, _ := regexp.MatchString(`\d{2}/\d{2}/\d{4} \d{2}:\d{2}:\d{2}`, data)
		if matched {
			dadosCompra.Data = data
			return
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		scrapeDate(c)
	}
}
func main() {
	url := "https://portalsped.fazenda.mg.gov.br/portalnfce/sistema/qrcode.xhtml?p=31250101928075004278650090004419571142667462%7C2%7C1%7C1%7C20422ea97778a2db22109f5c5b218e22fd62c05a"
	doc, err := fetch(url)
	if err != nil {
		fmt.Printf("Error fetching URL: %v\n", err)
		return
	}

	scrapeProducts(doc)
	scrapeSaleInfo(doc)
	scrapeStoreName(doc)
	scrapeStoreAddress(doc)
	scrapeDate(doc)
	dadosCompra.Produtos = produtos
	//scrapeTitles(doc)
}
