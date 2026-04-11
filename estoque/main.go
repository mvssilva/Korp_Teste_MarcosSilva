package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Produto struct {
	Codigo    string  `json:"codigo"`
	Descricao string  `json:"descricao"`
	Saldo     float64 `json:"saldo"`
}

var db *sql.DB

func conectarBanco() {
	// string padrao de conexao = usuario:senha@tcp(endereco:porta)/nome_do_banco
	dsn := "dev:senha123@tcp(127.0.0.1:3306)/korp_estoque"
	var err error

	// prepara conexao
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erro crítico ao preparar a conexao: ", err)
	}

	// testa se a senha e o banco estão certos
	if err := db.Ping(); err != nil {
		log.Fatal("Erro! O banco não responder: ", err)
	}

	fmt.Println("Conexão com MySQL estabelecida com sucesso!")
}

func criarTabela() {
	query := `
	CREATE TABLE IF NOT EXISTS produtos(
		codigo VARCHAR(50) PRIMARY KEY,
		descricao VARCHAR(255) NOT NULL,
		saldo DOUBLE NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Erro ao criar a tabela: ", err)
	}

	fmt.Println("Tabela 'produtos' pronta para uso!")

}

func cadastrarProduto(c *gin.Context) {
	var p Produto

	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Formato de dados invalido"})
		return
	}

	query := "INSERT INTO produtos (codigo, descricao, saldo) VALUES (?, ?, ?)"

	_, err := db.Exec(query, p.Codigo, p.Descricao, p.Saldo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao salvar no banco: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}

func listarProdutos(c *gin.Context) {
	linhas, err := db.Query("SELECT codigo, descricao, saldo FROM produtos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao buscar produtos no banco de dados!"})
		return
	}

	defer linhas.Close()

	var listaProdutos []Produto

	for linhas.Next() {
		var p Produto
		if err := linhas.Scan(&p.Codigo, &p.Descricao, &p.Saldo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao ler o produto"})
			return
		}
		listaProdutos = append(listaProdutos, p)
	}

	c.JSON(http.StatusOK, listaProdutos)
}

func baixarEstoque(c *gin.Context) {
	codigo := c.Param("codigo")

	var dados struct {
		Quantidade float64 `json:"quantidade"`
	}

	if err := c.ShouldBindJSON(&dados); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Quantidade invalida."})
		return
	}

	query := "UPDATE produtos SET saldo = saldo - ? WHERE codigo = ? AND saldo >= ?"
	resultado, err := db.Exec(query, dados.Quantidade, codigo, dados.Quantidade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro interno no banco."})
		return
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Produto nao encontrado ou saldo insuficiente"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensagem": "Estoque atualizado com sucesso."})
}

func main() {

	conectarBanco()
	criarTabela()

	router := gin.Default()

	router.POST("/produtos", cadastrarProduto)
	router.GET("/produtos", listarProdutos)
	router.PUT("/produtos/:codigo/baixar", baixarEstoque)

	router.Run(":8080")
}
