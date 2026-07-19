# evolution-go

Go client wrapper for Evolution with a focus on type safety, local payload validation, and contract testing.

## What this module is

This module encapsulates Evolution API REST routes into Go services:

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

It provides:

- Typed requests and responses
- Local validations before HTTP calls
- Standardized API error handling
- Service-level mocks for unit tests
- Integration tests with Docker environment

## Documentation and contract sources

This module is based on two sources:

1. Official Evolution API documentation

- https://doc.evolution-api.com/v2/

2. Official source code (routers, DTOs, and schemas)

- https://github.com/evolution-foundation/evolution-api
- https://github.com/evolution-foundation/evolution-api/blob/main/src/api/routes/sendMessage.router.ts
- https://github.com/evolution-foundation/evolution-api/blob/main/src/api/routes/group.router.ts
- https://github.com/evolution-foundation/evolution-api/blob/main/src/api/dto/sendMessage.dto.ts
- https://github.com/evolution-foundation/evolution-api/blob/main/src/api/dto/group.dto.ts

Important note:

- Some endpoints may change response shapes across versions/integrations.
- Current relevant example: fetchAllGroups returns a JSON array of groups, not an object with success.

## Instalação

Requirements:

- Go 1.26+

Command:

```bash
go get github.com/luiz-otavio/evolution-go/api
```

## How to use

Minimal example:

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
        InstanceName: "my-instance",
        Qrcode:       true,
        Integration:  "WHATSAPP-BAILEYS",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Created instance:", resp.Instance.InstanceName)
}
```

## Common flows

### 1) Connect an instance and send text

```go
ctx := context.Background()

_, _ = client.InstanceService().Create(ctx, evolution.InstanceCreateRequest{
    InstanceName: "my-instance",
    Qrcode:       true,
    Integration:  "WHATSAPP-BAILEYS",
})

conn, _ := client.InstanceService().Connect(ctx, "my-instance")
fmt.Println("QRCode (base64/data-uri):", conn.Base64)

msg, _ := client.MessageService().SendText(ctx, "my-instance", evolution.SendTextRequest{
    Number: "5511999999999",
    Text:   "Hello from evolution-go/api",
})
fmt.Println("Status:", msg.Status, "MsgID:", msg.Key.ID)
```

### 2) List groups and send a message to a group by JID

```go
groups, _ := client.GroupService().FetchAllGroups(ctx, "my-instance", false)

var privateJID string
for _, g := range groups {
    if g.Subject == "Private" && g.ID != "" {
        privateJID = g.ID
        break
    }
}

if privateJID != "" {
    _, _ = client.MessageService().SendText(ctx, "my-instance", evolution.SendTextRequest{
        Number: privateJID,
        Text:   "Test message to Private group",
    })
}
```

## Typed structures (summary)

- MessageResponse uses typed MessageKey (id, remoteJid, fromMe)
- ChatService.FindChats returns []ChatSummary
- GroupService.FetchAllGroups returns []GroupSummary
- InstanceCreateResponse and InstanceResponse use typed InstanceInfo

This reduces generic maps and improves autocomplete and compile-time safety.

## Erros

When the API returns an HTTP error, the module converts it to a structured error.

Example:

```go
_, err := client.InstanceService().Connect(ctx, "missing-instance")
if err != nil {
    fmt.Println(err)
}
```

## Testes

### Contract and unit tests

Run from the api directory:

```bash
go test ./... -count=1
```

### Integration tests

Dependencies:

- Working Docker environment
- Evolution API, Postgres, and Redis via testcontainers
- QR pairing when required

Useful environment variables (file [api/.env](api/.env)):

- EVOLUTION_BASE_URL
- EVOLUTION_API_KEY
- EVOLUTION_TEST_TARGET_NUMBER
- EVOLUTION_TEST_GROUP_JID
- EVOLUTION_TEST_GROUP_NAME

Run integration tests with verbose logs:

```bash
go test ./... -run Integration -v -count=1
```