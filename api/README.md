# evolution-go

Cliente Go (Wrapper) para Evolution com foco em segurança de tipos, validação local de payloads e testes de contrato.

## O que é este módulo

Este módulo encapsula as rotas REST da Evolution API em serviços Go:

- InstanceService
- MessageService
- ChatService
- GroupService
- LabelService
- SettingsService
- ProxyService
- WebhookService
- TemplateService
- CallService
- BusinessService

Ele oferece:

- Requests e responses tipados
- Validações locais antes da chamada HTTP
- Tratamento padronizado de erro da API
- Mocks por serviço para testes unitários
- Testes de integração com ambiente Docker

## Base de documentação e referência de contrato

Este módulo é baseado em duas fontes:

1. Documentação oficial Evolution API

- https://doc.evolution-api.com/v2/

2. Código-fonte oficial (router, dto e schemas)

- https://github.com/evolution-foundation/evolution-api
- https://github.com/evolution-foundation/evolution-api/blob/main/src/api/routes/sendMessage.router.ts
- https://github.com/evolution-foundation/evolution-api/blob/main/src/api/routes/group.router.ts
- https://github.com/evolution-foundation/evolution-api/blob/main/src/api/dto/sendMessage.dto.ts
- https://github.com/evolution-foundation/evolution-api/blob/main/src/api/dto/group.dto.ts

Observação importante:

- Alguns endpoints mudam formato de resposta entre versões e integrações.
- Exemplo atual relevante: fetchAllGroups retorna JSON array de grupos, não objeto com success.

## Instalação

Requisitos:

- Go 1.26+

Comando:

```bash
go get github.com/luiz-otavio/evolution-go/api
```

## Como usar

Exemplo mínimo:

```go
package main

import (
    "context"
    "fmt"
    "log"

    evolution "github.com/luiz-otavio/evolution-go/api"
)

func main() {
    client, err := evolution.NewEvolutionClient(evolution.EvolutionConfig{
        BaseURL: "http://localhost:8080",
        APIKey:  "test-api-key",
    })
    if err != nil {
        log.Fatal(err)
    }

    resp, err := client.InstanceService().Create(context.Background(), evolution.InstanceCreateRequest{
        InstanceName: "minha-instancia",
        Qrcode:       true,
        Integration:  "WHATSAPP-BAILEYS",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Instância criada:", resp.Instance.InstanceName)
}
```

## Fluxos comuns

### 1) Conectar instância e enviar texto

```go
ctx := context.Background()

_, _ = client.InstanceService().Create(ctx, evolution.InstanceCreateRequest{
    InstanceName: "minha-instancia",
    Qrcode:       true,
    Integration:  "WHATSAPP-BAILEYS",
})

conn, _ := client.InstanceService().Connect(ctx, "minha-instancia")
fmt.Println("QRCode (base64/data-uri):", conn.Base64)

msg, _ := client.MessageService().SendText(ctx, "minha-instancia", evolution.SendTextRequest{
    Number: "5511999999999",
    Text:   "Olá da evolution-go/api",
})
fmt.Println("Status:", msg.Status, "MsgID:", msg.Key.ID)
```

### 2) Listar grupos e enviar mensagem em grupo pelo JID

```go
groups, _ := client.GroupService().FetchAllGroups(ctx, "minha-instancia", false)

var privateJID string
for _, g := range groups {
    if g.Subject == "Private" && g.ID != "" {
        privateJID = g.ID
        break
    }
}

if privateJID != "" {
    _, _ = client.MessageService().SendText(ctx, "minha-instancia", evolution.SendTextRequest{
        Number: privateJID,
        Text:   "Mensagem de teste para o grupo Private",
    })
}
```

## Estrutura de tipos (resumo)

- MessageResponse usa MessageKey tipado (id, remoteJid, fromMe)
- ChatService.FindChats retorna []ChatSummary
- GroupService.FetchAllGroups retorna []GroupSummary
- InstanceCreateResponse e InstanceResponse usam InstanceInfo tipado

Isso reduz map genérico e melhora autocomplete e segurança de compilação.

## Erros

Quando a API retorna erro HTTP, o módulo converte para erro estruturado.

Exemplo:

```go
_, err := client.InstanceService().Connect(ctx, "instancia-inexistente")
if err != nil {
    fmt.Println(err)
}
```

## Testes

### Testes de contrato e unitários

Executar no diretório api:

```bash
go test ./... -count=1
```

### Testes de integração

Dependências:

- Docker funcional
- Evolution API, Postgres e Redis via testcontainers
- Pareamento de QR quando necessário

Variáveis de ambiente úteis (arquivo [api/.env](api/.env)):

- EVOLUTION_BASE_URL
- EVOLUTION_API_KEY
- EVOLUTION_TEST_TARGET_NUMBER
- EVOLUTION_TEST_GROUP_JID
- EVOLUTION_TEST_GROUP_NAME

Rodar integração com logs detalhados:

```bash
go test ./... -run Integration -v -count=1
```