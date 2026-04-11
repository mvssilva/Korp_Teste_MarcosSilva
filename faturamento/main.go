package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type ItemNota struct {
	ProdutoCodigo string  `json:"produto_codigo"`
	Quantidade    float64 `json:"quantidade"`
}

type NotaFiscal struct {
	ID     int        `json:"id"`
	Status string     `json:"status"`
	Itens  []ItemNota `json:"itens"`
}

var db *sql.DB

func conectarBanco() {

	dsn := "dev:senha123@tcp(127.0.0.1:3306)/korp_estoque"
	var err error

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erro critico ao preparar conexão: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Erro! O banco não respondeu: ", err)
	}

	fmt.Println("Conexão com MySQL (faturamento) estabelecida!")
}

func criarTabelaNotas() {
	queryNota := `
	CREATE TABLE IF NOT EXISTS notas_fiscais (
		id INT AUTO_INCREMENT PRIMARY KEY,
		status VARCHAR(20) DEFAULT 'Aberta'
	);`

	_, err := db.Exec(queryNota)
	if err != nil {
		log.Fatal("Erro ao criar a tabela de notas: ", err)
	}

	queryItens := `
	CREATE TABLE IF NOT EXISTS itens_nota (
		id INT AUTO_INCREMENT PRIMARY KEY,
		nota_id INT,
		produto_codigo VARCHAR(50),
		quantidade DOUBLE,
		FOREIGN KEY (nota_id) REFERENCES notas_fiscais(id)
	);`

	_, err = db.Exec(queryItens)
	if err != nil {
		log.Fatal("Erro ao criar tabela de itens: ", err)
	}

	fmt.Println("Tabelas de Faturamentos prontas para uso!")
}

func criarNota(c *gin.Context) {
	var nota NotaFiscal

	if err := c.ShouldBindJSON(&nota); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Formato inválido de nota."})
		return
	}

	resultado, err := db.Exec("INSERT INTO notas_fiscais (status) VALUES ('Aberta')")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao criar a nota principal."})
		return
	}

	notaID, _ := resultado.LastInsertId()
	nota.ID = int(notaID)
	nota.Status = "Aberta"

	for _, item := range nota.Itens {
		_, err := db.Exec("INSERT INTO itens_nota (nota_id, produto_codigo, quantidade) VALUES (?, ?, ?)", nota.ID, item.ProdutoCodigo, item.Quantidade)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao salvar item da nota."})
			return
		}
	}

	c.JSON(http.StatusCreated, nota)
}

func imprimirNota(c *gin.Context) {
	notaID := c.Param("id")

	var status string
	err := db.QueryRow("SELECT status FROM notas_fiscais WHERE id = ?", notaID).Scan(&status)
	if err != nil || status != "Aberta" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Nota não encontrada ou já foi fechada."})
		return
	}
	linhas, err := db.Query("SELECT produto_codigo, quantidade FROM itens_nota WHERE nota_id = ?", notaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao buscar itens da nota"})
	}
	defer linhas.Close()

	for linhas.Next() {
		var codigo string
		var qtd float64
		linhas.Scan(&codigo, &qtd)

		corpoRequisicao := fmt.Sprintf(`{"quantidade": %f}`, qtd)

		urlEstoque := "http://localhost:8080/produtos/" + codigo + "/baixar"

		req, _ := http.NewRequest("PUT", urlEstoque, bytes.NewBuffer([]byte(corpoRequisicao)))
		req.Header.Set("Content-Type", "application/json")

		cliente := &http.Client{}
		resposta, _ := cliente.Do(req)

		if resposta.StatusCode != http.StatusOK {
			c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha no estoque para o produto " + codigo})
			return
		}
	}

	db.Exec("UPDATE notas_fiscais SET status = 'Fechada' WHERE id = ?", notaID)

	c.JSON(http.StatusOK, gin.H{"mensagem": "Nota impressa e fechado com sucesso, o saldo foi descontado."})
}

func main() {
	conectarBanco()
	criarTabelaNotas()

	router := gin.Default()
	router.POST("/notas", criarNota)

	router.POST("/notas/:id/imprimir", imprimirNota)
	router.Run(":8081")
}
