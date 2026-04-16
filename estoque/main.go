package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Representa a entidade Produto no banco de dados
type Produto struct {
	Codigo    string  `json:"codigo"`
	Descricao string  `json:"descricao"`
	Saldo     float64 `json:"saldo"`
}

var db *sql.DB

func conectarBanco() {
	// Configura a string de conexão (DSN) com o MySQL
	dsn := "dev:senha123@tcp(127.0.0.1:3306)/korp_estoque"
	var err error

	// Inicializa o pool de conexões
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erro crítico ao preparar a conexão: ", err)
	}

	// Verifica a disponibilidade e autenticação do banco de dados
	if err := db.Ping(); err != nil {
		log.Fatal("Erro! O banco não está respondendo: ", err)
	}

	fmt.Println("Conexão com MySQL estabelecida com sucesso!")
}

func criarTabela() {
	// Garante a existência da tabela principal de produtos na inicialização
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
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Formato de dados inválido"})
		return
	}

	// 1. Remove espaços vazios do começo e do final
	p.Codigo = strings.TrimSpace(p.Codigo)
	p.Descricao = strings.TrimSpace(p.Descricao)

	// 2. Padroniza tudo para MAIÚSCULO (Padrão de ERP e Notas Fiscais)
	p.Codigo = strings.ToUpper(p.Codigo)
	p.Descricao = strings.ToUpper(p.Descricao)

	// Persiste o novo produto no banco de dados
	query := "INSERT INTO produtos (codigo, descricao, saldo) VALUES (?, ?, ?)"
	_, err := db.Exec(query, p.Codigo, p.Descricao, p.Saldo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao salvar no banco: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}

func listarProdutos(c *gin.Context) {
	// Consulta todos os produtos cadastrados
	linhas, err := db.Query("SELECT codigo, descricao, saldo FROM produtos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao buscar produtos no banco de dados!"})
		return
	}
	defer linhas.Close()

	var listaProdutos []Produto

	// Itera sobre os resultados e popula a fatia (slice) de produtos
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

	// Extrai a quantidade a ser deduzida do payload
	if err := c.ShouldBindJSON(&dados); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Quantidade inválida."})
		return
	}

	// Executa a baixa de forma atômica, validando se há saldo suficiente na mesma query
	query := "UPDATE produtos SET saldo = saldo - ? WHERE codigo = ? AND saldo >= ?"
	resultado, err := db.Exec(query, dados.Quantidade, codigo, dados.Quantidade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro interno no banco."})
		return
	}

	// Verifica se a query realmente alterou alguma linha (previne saldo negativo)
	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Produto não encontrado ou saldo insuficiente"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensagem": "Estoque atualizado com sucesso."})
}

// Middleware para habilitar CORS (Cross-Origin Resource Sharing) com o Angular
func controleAcesso(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Intercepta e aprova requisições de preflight do navegador
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}

// Busca um produto específico pelo código para validação no frontend
func buscarProduto(c *gin.Context) {
	codigo := c.Param("codigo")
	var p Produto

	query := "SELECT codigo, descricao, saldo FROM produtos WHERE codigo = ?"
	err := db.QueryRow(query, codigo).Scan(&p.Codigo, &p.Descricao, &p.Saldo)

	// Se o banco não encontrar a linha, devolve Erro 404
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Produto não encontrado"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao consultar banco de dados"})
		return
	}

	c.JSON(http.StatusOK, p)
}

func main() {
	// Inicialização da infraestrutura
	conectarBanco()
	criarTabela()

	router := gin.Default()

	// Configuração de rotas e middlewares
	router.Use(controleAcesso)

	router.POST("/produtos", cadastrarProduto)
	router.GET("/produtos", listarProdutos)
	router.GET("/produtos/:codigo", buscarProduto)
	router.PUT("/produtos/:codigo/baixar", baixarEstoque)

	// Inicializa o servidor HTTP na porta 8080
	router.Run(":8080")
}
