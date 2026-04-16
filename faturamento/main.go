package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv" //

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Estrutura para receber a pergunta do Angular
type RequisicaoChat struct {
	Mensagem string `json:"mensagem"`
}

// Representa o item individual vinculado a uma nota fiscal
type ItemNota struct {
	ProdutoCodigo string  `json:"produto_codigo"`
	Quantidade    float64 `json:"quantidade"`
}

// Representa a entidade principal da Nota Fiscal
type NotaFiscal struct {
	ID     int        `json:"id"`
	Status string     `json:"status"`
	Itens  []ItemNota `json:"itens"`
}

var db *sql.DB

func conectarBanco() {
	// Configura a string de conexão (DSN) com o MySQL
	dsn := "dev:senha123@tcp(127.0.0.1:3306)/korp_estoque"
	var err error

	// Inicializa o pool de conexões
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erro crítico ao preparar conexão: ", err)
	}

	// Verifica a disponibilidade do banco de dados
	if err := db.Ping(); err != nil {
		log.Fatal("Erro! O banco não respondeu: ", err)
	}

	fmt.Println("Conexão com MySQL (faturamento) estabelecida!")
}

func criarTabelaNotas() {
	// Garante a existência da tabela principal (Cabeçalho da Nota)
	queryNota := `
    CREATE TABLE IF NOT EXISTS notas_fiscais (
        id INT AUTO_INCREMENT PRIMARY KEY,
        status VARCHAR(20) DEFAULT 'Aberta'
    );`

	_, err := db.Exec(queryNota)
	if err != nil {
		log.Fatal("Erro ao criar a tabela de notas: ", err)
	}

	// Garante a existência da tabela de itens (Relacionamento 1:N)
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

	// Valida o payload recebido do frontend
	if err := c.ShouldBindJSON(&nota); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Formato inválido de nota."})
		return
	}

	// 1. Cria o cabeçalho da nota com status inicial "Aberta"
	resultado, err := db.Exec("INSERT INTO notas_fiscais (status) VALUES ('Aberta')")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao criar a nota principal."})
		return
	}

	// Recupera o ID gerado pelo banco para vincular aos itens
	notaID, _ := resultado.LastInsertId()
	nota.ID = int(notaID)
	nota.Status = "Aberta"

	// 2. Persiste cada item vinculado ao ID da nota principal
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

	// Impede a reemissão de notas já processadas/fechadas
	var status string
	err := db.QueryRow("SELECT status FROM notas_fiscais WHERE id = ?", notaID).Scan(&status)
	if err != nil || status != "Aberta" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Nota não encontrada ou já foi fechada."})
		return
	}

	// Busca todos os itens da nota para processar a baixa no estoque
	linhas, err := db.Query("SELECT produto_codigo, quantidade FROM itens_nota WHERE nota_id = ?", notaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao buscar itens da nota"})
		return
	}
	defer linhas.Close()

	// Itera sobre os itens e notifica o microsserviço de Estoque
	for linhas.Next() {
		var codigo string
		var qtd float64
		linhas.Scan(&codigo, &qtd)

		// Estrutura o payload e a URL para a API de Estoque
		corpoRequisicao := fmt.Sprintf(`{"quantidade": %d}`, int(qtd))
		urlEstoque := "http://localhost:8080/produtos/" + codigo + "/baixar"

		req, errReq := http.NewRequest("PUT", urlEstoque, bytes.NewBuffer([]byte(corpoRequisicao)))
		if errReq != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao montar requisição para o estoque"})
			return
		}

		req.Header.Set("Content-Type", "application/json")

		// Executa a chamada HTTP para o microsserviço
		cliente := &http.Client{}
		resposta, errResp := cliente.Do(req)

		// Trata falha severa de rede (Ex: Estoque offline)
		if errResp != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"erro": "Falha de comunicação com o microsserviço de Estoque."})
			return
		}
		defer resposta.Body.Close()

		// Trata recusa do estoque (Ex: Saldo insuficiente ou produto não encontrado)
		if resposta.StatusCode != http.StatusOK {
			c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha no estoque para o produto " + codigo})
			return
		}
	}

	// Marca a nota como "Fechada" após abater todos os saldos com sucesso
	db.Exec("UPDATE notas_fiscais SET status = 'Fechada' WHERE id = ?", notaID)

	c.JSON(http.StatusOK, gin.H{"mensagem": "Nota impressa e fechada com sucesso, o saldo foi descontado."})
}

func controleAcesso(c *gin.Context) {
	// Habilita CORS para comunicação com o Angular
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

// Rota inteligente que processa a pergunta e devolve a resposta da IA
func chatComGemini(c *gin.Context) {
	var req RequisicaoChat

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Mensagem inválida."})
		return
	}

	// SEGURANÇA: Lendo a chave da variável de ambiente do Windows/Linux
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "API Key não configurada no servidor."})
		return
	}

	// Inicializando o cliente do Gemini
	ctx := context.Background()
	cliente, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao conectar com a IA."})
		return
	}
	defer cliente.Close()

	// Usando o modelo flash, que é o mais rápido para chatbots
	modelo := cliente.GenerativeModel("gemini-2.5-flash")

	// System Prompt Avançado: O Cérebro do seu ERP
	// System Prompt Avançado com Blindagem de Dados (Data Leakage Prevention)
	promptSistema := `Você é o KorpAssist, o assistente virtual integrado ao sistema ERP Korp.
    
    [CONTEXTO TÉCNICO INTERNO - APENAS PARA SEU RACIOCÍNIO]
    - O sistema é dividido em Microsserviços.
    - O Faturamento consome a API do Estoque automaticamente para baixar produtos.
    - O saldo de um produto nunca fica negativo no banco de dados.
    - Produtos nunca são deletados, apenas marcados caso saiam de linha.
    
    [DIRETRIZES DE COMUNICAÇÃO - REGRA DE OURO]
    1. O seu público-alvo é o USUÁRIO FINAL do sistema (faturistas, estoquistas, gerentes).
    2. VOCÊ ESTÁ ESTRITAMENTE PROIBIDO de mencionar detalhes de infraestrutura na sua resposta, 
	como: portas de servidor (ex: 8080, 8081), linguagens de programação (Go, Angular), rotas de API, ou estrutura de banco de dados.
    3. Explique os fluxos sempre pela perspectiva do negócio e do uso da tela.
    4. Responda de forma profissional, direta e em no máximo 1 parágrafo.
    
    Pergunta do usuário: `

	// Enviando a pergunta
	resposta, err := modelo.GenerateContent(ctx, genai.Text(promptSistema+req.Mensagem))

	// 🔥 CÓDIGO DE DEBUG ADICIONADO AQUI
	if err != nil {
		fmt.Println("\n=== ERRO DETALHADO DA API DO GOOGLE ===")
		fmt.Println(err)
		fmt.Println("========================================\n")

		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha na API: " + err.Error()})
		return
	}
	// Extraindo e limpando o texto da resposta que o Gemini devolveu
	var textoResposta string
	if len(resposta.Candidates) > 0 && len(resposta.Candidates[0].Content.Parts) > 0 {
		if texto, ok := resposta.Candidates[0].Content.Parts[0].(genai.Text); ok {
			textoResposta = string(texto)
		}
	}

	c.JSON(http.StatusOK, gin.H{"resposta": textoResposta})
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Arquivo .env não encontrado. O sistema tentará usar as variáveis do Windows/Linux.")
	}

	// Inicialização da infraestrutura
	conectarBanco()
	criarTabelaNotas()

	router := gin.Default()

	// Configuração de rotas e middlewares
	router.Use(controleAcesso)
	router.POST("/notas", criarNota)
	router.POST("/notas/:id/imprimir", imprimirNota)
	router.POST("/chat", chatComGemini)

	// Inicializa o servidor HTTP na porta 8081
	router.Run(":8081")
}
