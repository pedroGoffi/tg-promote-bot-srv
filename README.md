# Bot-Manager

Bot-Manager é um servidor escrito em Go que gerencia um bot para download de imagens.

## Características
- Desenvolvido em Go utilizando `gin-gonic` para gerenciamento de rotas.
- Suporte para download de imagens a partir de links do Telegram.
- Armazenamento local de imagens no diretório `./images`.
- Servidor HTTP para disponibilizar imagens armazenadas.
- Upload de imagens para um servidor externo.

## Instalação e Configuração

1. Clone este repositório:
   ```sh
   git clone https://github.com/seu-usuario/bot-manager.git
   cd bot-manager
   ```
2. Instale as dependências:
   ```sh
   go mod tidy
   ```
3. Configure as variáveis de ambiente no arquivo `.env` (caso necessário).

## Uso

### Iniciando o Servidor

Para iniciar o servidor, execute:
```sh
go run main.go
```
O servidor rodará na porta configurada nas variáveis de ambiente.

### Rotas Disponíveis

#### `POST /api/tgimg`
Baixa uma imagem a partir de um link do Telegram.

**Parâmetros:**
```json
{
  "telegram_link": "URL_da_imagem_no_Telegram"
}
```

**Resposta de Sucesso:**
```json
{
  "message": "Imagem baixada com sucesso",
  "endpoint": "/static/nome_da_imagem"
}
```

### Servindo Imagens
As imagens são armazenadas no diretório `./images` e podem ser acessadas pelo endpoint `/static/{nome_da_imagem}`.

## Estrutura do Projeto

```
 bot-manager/
 ├── internal/
 │   ├── config/        # Configuração do servidor
 │   ├── handlers/      # Handlers para as requisições HTTP
 │   ├── routes/        # Configuração das rotas
 │   ├── uploader/      # Módulo para upload de imagens
 ├── images/           # Diretório para armazenar imagens baixadas
 ├── main.go           # Arquivo principal
 ├── go.mod            # Dependências do projeto
 ├── README.md         # Documentação do projeto
```

## Autor
Pedro Henrique Goffi de Paulo

