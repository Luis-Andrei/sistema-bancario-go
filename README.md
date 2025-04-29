# Sistema Bancário em Go

Um servidor HTTP em Go que gerencia operações bancárias simples, utilizando PostgreSQL para armazenamento de dados e proteção contra ataques CSRF.

## Funcionalidades

- Criação de contas bancárias
- Consulta de saldo
- Depósitos
- Saques (com verificação de saldo)
- Proteção CSRF

## Requisitos

- Go 1.21 ou superior
- PostgreSQL
- Git

## Configuração

1. Clone o repositório:
```bash
git clone https://github.com/seu-usuario/nome-do-repositorio.git
cd nome-do-repositorio
```

2. Instale as dependências:
```bash
go mod tidy
```

3. Configure o banco de dados:
- Instale o PostgreSQL
- Crie um banco de dados chamado "banco"
- Atualize as credenciais em `config/database.go`

4. Execute o servidor:
```bash
go run main.go
```

## Endpoints

- `POST /accounts` - Criar nova conta
- `GET /accounts/{id}` - Consultar conta
- `POST /accounts/{id}/deposit` - Realizar depósito
- `POST /accounts/{id}/withdraw` - Realizar saque

## Exemplos de Uso

```bash
# Criar conta
curl -X POST http://localhost:8080/accounts -H "Content-Type: application/json" -d '{"number":"12345","balance":1000}'

# Consultar conta
curl http://localhost:8080/accounts/1

# Realizar depósito
curl -X POST http://localhost:8080/accounts/1/deposit -H "Content-Type: application/json" -d '{"amount":500}'

# Realizar saque
curl -X POST http://localhost:8080/accounts/1/withdraw -H "Content-Type: application/json" -d '{"amount":200}'
```

## Estrutura do Projeto

```
.
├── config/
│   └── database.go    # Configurações do banco de dados
├── models/
│   └── account.go     # Modelo de conta bancária
├── handlers/
│   └── account_handler.go  # Handlers HTTP
├── main.go            # Ponto de entrada da aplicação
├── go.mod             # Dependências do projeto
└── README.md          # Documentação
```

## Contribuição

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -m 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo LICENSE para detalhes. 