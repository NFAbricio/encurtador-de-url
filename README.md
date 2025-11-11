# Encurtador de URL em Go

Este é um projeto simples de um encurtador de URL desenvolvido em Go. A aplicação permite encurtar URLs longas para um formato com 8 caracteres e redirecionar para a URL original quando o link encurtado é acessado.

## Funcionalidades

- Encurta URLs longas para um formato mais amigável.
- Redireciona para a URL original ao acessar o link encurtado.
- API RESTful para integração com outras aplicações.
- Armazenamento em memória (para fins de demonstração).

## Endpoints da API

- `POST /api/shorten`: Encurta uma nova URL.
  - **Body (JSON):** `{"url": "sua-url-longa"}`
  - **Resposta de Sucesso (201 Created):** `{"data": "codigo-curto"}`
- `GET /{code}`: Redireciona para a URL original.

## Como Executar o Projeto

### Pré-requisitos

- [Go](https://golang.org/dl/) instalado (versão 1.22.3 ou superior).

### Passos

1. **Clone o repositório:**
   ```bash
   git clone https://github.com/seu-usuario/encurtador-go.git
   cd encurtador-go
   ```

2. **Instale as dependências:**
   ```bash
   go mod tidy
   ```

3. **Execute a aplicação:**
   ```bash
   go run main.go
   ```

A aplicação estará disponível em `http://localhost:8080`.

## Tecnologias Utilizadas

- **Go:** Linguagem de programação principal.
- **Gorilla Mux:** Um poderoso roteador de URL e dispatcher para Go, usado para gerenciar as rotas da API.
- **Gorilla Handlers:** Uma coleção de middlewares para servidores HTTP, utilizado para logging e recuperação de panics.
- **xid:** Para geração de IDs únicos.