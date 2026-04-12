# Korp_Teste_MarcosSilva
Este projeto foi desenvolvido como parte do desafio técnico para a vaga de Desenvolvedor na Korp/Viasoft. O sistema consiste em uma aplicação web fullstack para o cadastro de produtos e gestão de notas fiscais, operando sob uma arquitetura de microsserviços.

## Apresentação
- Vídeo Demonstrativo: [Link para o vídeo (Google Drive/OneDrive)] 
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
    
## Como Executar o Projeto
Pré-requisitos
- Go 1.20+
- Node.js & Angular CLI
- MySQL Server

Passo a Passo
  1. Banco de Dados: Execute os scripts SQL localizados em /database/init.sql.
  2. Microsserviço de Estoque:
  ```bash
    cd backend-estoque
    go run main.go
  ```
  3. Microsserviço de Faturamento:
  ```bash
    cd backend-faturamento
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

## Fotos do estágio atual:
<img width="795" height="781" alt="image" src="https://github.com/user-attachments/assets/ec8ec716-9d9a-416a-9aaa-b25f548b3255" />

Desenvolvido por Marcos Silva, Obrigado!
