# evolution-go

SDK em Go para integração com a Evolution API.

Este repositório organiza os módulos do projeto por responsabilidades (cliente API, testes, scripts auxiliares), com foco em:

- simplicidade de uso para integrações WhatsApp
- tipagem forte dos contratos HTTP
- testes de contrato e integração
- evolução incremental baseada no contrato upstream da Evolution API

## Visão geral

A ideia do projeto é fornecer uma camada Go para consumir os endpoints da Evolution API com validação local de payload, responses tipadas e mocks para testes.

## Módulos do workspace

### `api/`

Módulo principal do SDK.

Contém:

- cliente HTTP e configuração (`NewEvolutionClient`)
- serviços por domínio (instance, message, chat, group, settings, etc.)
- modelos de request/response
- mocks para testes unitários
- suíte de testes (contrato + integração)

Documentação completa do módulo: [api/README.md](api/README.md)

### `scripts/`

Scripts utilitários de suporte ao desenvolvimento.

Atualmente:

- `find-affected-services.sh`: ajuda a identificar serviços impactados por mudanças

## Estrutura do workspace

- `go.work`: define os módulos Go ativos no workspace
- `api/go.mod`: dependências e versão Go do módulo `api`

## Começando rápido

Pré-requisitos:

- Go 1.26+

Instalar dependências do módulo API:

```bash
cd api
go mod tidy
```

Executar testes de contrato/unitários:

```bash
go test ./... -count=1
```

## Testes de integração

Os testes de integração do módulo `api` usam Docker (testcontainers) e variáveis de ambiente.

Consulte:

- [api/README.md](api/README.md)
- [api/.env.example](api/.env.example)

Execução típica:

```bash
cd api
go test ./... -run Integration -v -count=1
```

## Referências

- Documentação Evolution API: https://doc.evolution-api.com/v2/
- Repositório oficial Evolution: https://github.com/evolution-foundation