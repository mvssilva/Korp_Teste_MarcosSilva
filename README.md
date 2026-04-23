# SISTEMA DE CONTROLE E EMISSOES DE NOTA
Este projeto foi desenvolvido como parte do desafio técnico para a vaga de Desenvolvedor na Korp/Viasoft. O sistema consiste em uma aplicação web fullstack para o cadastro de produtos e gestão de notas fiscais, operando sob uma arquitetura de microsserviços.

## Apresentação
- Vídeo Demonstrativo: [[Link para o vídeo (Google Drive)] ](https://drive.google.com/file/d/1PyCXHyD3WrtSA1CVaJsRJQgncP_8qc5z/view?usp=drive_link)
    - O vídeo contém a demonstração das telas, funcionalidades e o detalhamento técnico da solução.

## Arquitetura do Sistema
- O sistema foi desenhado seguindo o requisito obrigatório de Microsserviços:
  -   Serviço de Estoque (Porta 8080): Responsável pelo controle de produtos e saldos. Desenvolvido em Go utilizando o framework Gin.
  -   Serviço de Faturamento (Porta 8081): Gestão de notas fiscais e processamento de baixas. Desenvolvido em Go com Gin.
  -   Frontend (Angular): Interface Single Page Application (SPA) para interação com o usuário.
- Banco de Dados
  - MySQL: Utilizado para a persistência real dos dados (Produtos e Notas Fiscais).

## Tecnologias e Bibliotecas Utilizadas
- Frontend (Angular)
  - Ciclos de Vida: Utilização de ngOnInit para inicialização de dados e consumo de APIs.
  - RxJS: Uso de Observable e subscribe nos serviços para gerenciamento de chamadas assíncronas HTTP.
  - Componentes Visuais: Desenvolvimento de componentes customizados com CSS (seguindo a identidade visual da Korp).
  - Bibliotecas: @angular/router para gerenciamento de rotas SPA.
- Backend (Go)
  - Framework: Gin Gonic para roteamento HTTP e gerenciamento de Middleware (CORS).
  - Gerenciamento de Dependências: Go Modules (go.mod).
  - Tratamento de Erros: Implementado através de validações de entrada e respostas com códigos de status HTTP apropriados (400, 404, 500) para feedback ao usuário.
- Inteligência Artificial
    - Google Gemini API: Integrada ao microsserviço de faturamento para suporte ao usuário.
    - Generative AI SDK: Utilizado para processamento de linguagem natural e respostas contextualizadas.
    - Dotenv (godotenv): Gerenciamento seguro de chaves de API e variáveis de ambiente.

## Funcionalidades Implementadas
- Conforme o escopo do projeto:
  1. Gestão de Estoque: Cadastro de produtos com código, descrição e saldo.
  2. Gestão de Notas Fiscais:
     - Criação de notas com numeração sequencial e status "Aberta".
     - Inclusão de múltiplos produtos por nota.
  3. Impressão e Baixa:
     - Botão intuitivo de impressão com indicador de processamento.
     - Atualização automática do status da nota para "Fechada" após a emissão.
     - Atualização de Saldo: O sistema realiza o cálculo e abate o saldo dos produtos no estoque conforme a quantidade utilizada na nota.
  4. KorpAssist (IA):
        - Chatbot integrado para sanar dúvidas sobre o fluxo de estoque e faturamento.
        - Engenharia de Prompt: Blindagem contra vazamento de dados técnicos (Data Leakage) e foco em regras de negócio.
          
## Como Executar o Projeto
Pré-requisitos
- Go 1.20+
- Node.js & Angular CLI
- MySQL Server
- Gemini API Key (Google AI Studio)

Passo a Passo
  1. Banco de Dados: Execute os scripts SQL localizados em /database/init.sql.
  2. Configuração de Variáveis de Ambiente:
    - No diretório do serviço de faturamento, crie um arquivo `.env`.
    - Adicione a linha: `GEMINI_API_KEY=sua_chave_aqui`.
  3. Microsserviço de Estoque:
  ```bash
    cd backend-estoque
    go mod tidy
    go run main.go
  ```
  3. Microsserviço de Faturamento:
  ```bash
    cd backend-faturamento
    go mod tidy
    go run main.go
  ```
  4. Frontend:
  ```bash
    cd frontend
    npm install
    ng serve
  ```
  5. Acesse http://localhost:4200 no seu navegador.

## Tratamento de Falhas
Implementado feedback visual (Toast Notifications) para alertar o usuário caso um dos microsserviços esteja offline ou ocorra um erro de processamento, garantindo a recuperação e transparência do sistema.

Desenvolvido por Marcos Silva, Obrigado!
