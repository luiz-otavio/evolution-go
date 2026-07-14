# evolution

Cliente Go para a [Evolution API](https://docs.evolutionfoundation.com.br/evolution-api) — API REST completa para WhatsApp com suporte multi-provedor.

## Instalação

```bash
go get github.com/luiz-otavio/evolution-go/api
```

## Início rápido

```go
import evolution "github.com/luiz-otavio/evolution-go/api"

client, err := evolution.NewEvolutionClient(evolution.EvolutionConfig{
    BaseURL: evolution.BaseURLLocal, // ou BaseURLProduction, ou uma URL customizada
    APIKey:  "sua-api-key",          // global ou por instância
})
if err != nil {
    log.Fatal(err)
}
```

### Ambientes disponíveis

| Constante            | URL                              |
|----------------------|----------------------------------|
| `BaseURLLocal`       | `http://localhost:8080`          |
| `BaseURLProduction`  | `https://api.evolution-api.com`  |

A autenticação é feita via header `apikey` (global ou específico da instância).

---

## Serviços

O cliente expõe os seguintes serviços:

```go
client.InstanceService() // Instâncias (criar, conectar, status, logout...)
client.MessageService()  // Envio de mensagens (texto, mídia, botões, lista...)
client.ChatService()     // Chats, contatos, mensagens e perfil
client.GroupService()    // Grupos e participantes
client.LabelService()    // Etiquetas
client.SettingsService() // Configurações da instância
client.ProxyService()    // Proxy
client.WebhookService()  // Webhook e WebSocket
client.TemplateService() // Templates do WhatsApp Business
client.CallService()     // Chamadas
client.BusinessService() // Catálogo e coleções (WhatsApp Business)
```

---

## Instância

```go
ctx := context.Background()

// Criar instância
created, err := client.InstanceService().Create(ctx, evolution.InstanceCreateRequest{
    InstanceName: "minha-instancia",
    Qrcode:       true,
    Integration:  "WHATSAPP-BAILEYS",
})

// Conectar (retorna QR code / pairing code)
conn, err := client.InstanceService().Connect(ctx, "minha-instancia")
fmt.Println("QR base64:", conn.Base64)

// Estado da conexão
state, err := client.InstanceService().ConnectionState(ctx, "minha-instancia")
fmt.Println("Estado:", state.Instance.State)

// Listar todas
instances, err := client.InstanceService().FetchInstances(ctx)

// Reiniciar / logout / deletar
_, err = client.InstanceService().Restart(ctx, "minha-instancia")
_, err = client.InstanceService().Logout(ctx, "minha-instancia")
_, err = client.InstanceService().Delete(ctx, "minha-instancia")

// Presença
_, err = client.InstanceService().SetPresence(ctx, "minha-instancia",
    evolution.SetPresenceRequest{Presence: evolution.PresenceAvailable})
```

---

## Mensagens

```go
// Texto
_, err := client.MessageService().SendText(ctx, "minha-instancia",
    evolution.SendTextRequest{
        Number:      "5511999999999",
        TextMessage: evolution.TextMessage{Text: "Olá!"},
    })

// Mídia (Media é URL ou base64)
_, err = client.MessageService().SendMedia(ctx, "minha-instancia",
    evolution.SendMediaRequest{
        Number:    "5511999999999",
        MediaType: evolution.MediaTypeImage,
        Media:     "https://exemplo.com/imagem.png",
        Caption:   "Legenda",
    })

// Localização
_, err = client.MessageService().SendLocation(ctx, "minha-instancia",
    evolution.SendLocationRequest{
        Number:    "5511999999999",
        Latitude:  -23.55,
        Longitude: -46.63,
        Name:      "São Paulo",
    })

// Enquete
_, err = client.MessageService().SendPoll(ctx, "minha-instancia",
    evolution.SendPollRequest{
        Number:          "5511999999999",
        Name:            "Qual sua cor favorita?",
        SelectableCount: 1,
        Values:          []string{"Azul", "Verde", "Vermelho"},
    })
```

Também disponíveis: `SendButtons`, `SendList`, `SendContact`, `SendReaction`, `SendTemplate`.

---

## Chat

```go
// Verificar números no WhatsApp
res, err := client.ChatService().CheckWhatsAppNumbers(ctx, "minha-instancia",
    evolution.WhatsAppNumbersRequest{Numbers: []string{"5511999999999"}})

// Buscar contatos
contacts, err := client.ChatService().FindContacts(ctx, "minha-instancia",
    evolution.Query{Take: 20})

// Buscar mensagens
msgs, err := client.ChatService().FindMessages(ctx, "minha-instancia",
    evolution.Query{Take: 50})

// Arquivar chat / marcar como lido
_, err = client.ChatService().ArchiveChat(ctx, "minha-instancia",
    evolution.ArchiveChatRequest{Number: "5511999999999", Archive: true})

// Perfil
_, err = client.ChatService().UpdateProfileName(ctx, "minha-instancia",
    evolution.UpdateProfileNameRequest{Name: "Meu Nome"})
```

---

## Grupos

```go
// Criar grupo
grp, err := client.GroupService().Create(ctx, "minha-instancia",
    evolution.CreateGroupRequest{
        Subject:      "Meu Grupo",
        Participants: []string{"5511999999999"},
    })

// Info e participantes
info, err := client.GroupService().FindGroupInfos(ctx, "minha-instancia", "1203...@g.us")
parts, err := client.GroupService().Participants(ctx, "minha-instancia", "1203...@g.us")

// Adicionar / remover / promover / rebaixar
_, err = client.GroupService().UpdateParticipant(ctx, "minha-instancia",
    evolution.UpdateParticipantRequest{
        GroupJid:     "1203...@g.us",
        Action:       evolution.ParticipantActionAdd,
        Participants: []string{"5511888888888"},
    })
```

---

## Etiquetas, Configurações, Proxy e Webhook

```go
labels, err := client.LabelService().FindLabels(ctx, "minha-instancia")
_, err = client.LabelService().HandleLabel(ctx, "minha-instancia",
    evolution.HandleLabelRequest{Name: "urgente", Type: evolution.LabelTypeChat, ID: "123", Action: evolution.LabelActionAdd})

settings, err := client.SettingsService().Find(ctx, "minha-instancia")
_, err = client.SettingsService().Set(ctx, "minha-instancia",
    evolution.Settings{RejectCall: true, AlwaysOnline: true})

proxy, err := client.ProxyService().Find(ctx, "minha-instancia") // *Proxy (nil se não configurado)
_, err = client.ProxyService().Set(ctx, "minha-instancia",
    evolution.Proxy{Enabled: true, ProxyHost: "1.2.3.4", ProxyPort: "8080", ProxyProtocol: evolution.ProxyProtocolHTTP})

webhook, err := client.WebhookService().FindWebhook(ctx, "minha-instancia") // *Webhook
_, err = client.WebhookService().SetWebhook(ctx, "minha-instancia",
    evolution.SetWebhookRequest{Enabled: true, URL: "https://meu-servidor.com/webhook", Events: []string{"MESSAGES_UPSERT"}})
```

---

## Templates, Chamadas e Business

```go
tmpl, err := client.TemplateService().Create(ctx, "minha-instancia",
    evolution.CreateTemplateRequest{
        Name:       "boas_vindas",
        Category:   evolution.TemplateCategoryMarketing,
        Language:   "pt_BR",
        Components: []map[string]any{{"type": "BODY", "text": "Olá {{1}}"}},
    })

_, err = client.CallService().Offer(ctx, "minha-instancia",
    evolution.OfferCallRequest{Number: "5511999999999", Offer: map[string]any{}})

catalog, err := client.BusinessService().GetCatalog(ctx, "minha-instancia",
    evolution.BusinessNumberRequest{Number: "5511999999999"})
```

---

## Erros

Respostas de erro da Evolution API seguem o formato `{ "success": false, "error": { "code", "message" }, "meta": {...} }`.
Erros são retornados como `*evolution.EvolutionError`:

```go
_, err := client.InstanceService().Connect(ctx, "inexistente")
var apiErr *evolution.EvolutionError
if errors.As(err, &apiErr) {
    fmt.Println(apiErr.Code, apiErr.Message)
}
```

---

## Mocks

Cada serviço possui um mock (`Mock<Serviço>Service`) com campos de função (`...Fn`) para uso em testes:

```go
mock := &evolution.MockInstanceService{
    ConnectFn: func(ctx context.Context, instanceName string) (evolution.ConnectInstanceResponse, error) {
        return evolution.ConnectInstanceResponse{Code: "2@fake"}, nil
    },
}
```

> Nota: `SendMedia` e `UpdateProfilePicture` são implementados via corpo JSON (URL/base64), enquanto a especificação OpenAPI também documenta upload `multipart/form-data`.
