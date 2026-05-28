# eco

CLI em Go para criar e gerenciar projetos de API REST em Go.

## Instalação

```bash
go install github.com/dixavier27/eco/cmd/eco@latest
```

Requer Go 1.26 ou superior. Verifique com `eco doctor`.

## Comandos

| Comando | Descrição |
|---|---|
| `eco version` | Imprime a versão do binário. |
| `eco doctor` | Verifica se a toolchain Go está instalada e na versão mínima. |
| `eco new <nome>` | Cria um novo projeto de API REST em `./<nome>`. |

Cada comando aceita `--help`.

## Criando um projeto

```bash
eco new minha-api --module github.com/voce/minha-api
cd minha-api
go run ./cmd/api
# em outro terminal:
curl http://localhost:8080/healthz
```

Estrutura gerada:

```
minha-api/
├── cmd/api/main.go          # entrypoint
├── internal/
│   ├── http/                # servidor e rotas
│   └── config/              # configuração via env vars
├── api/openapi.yaml         # contrato OpenAPI 3 (stub)
├── configs/config.yaml      # configuração local (stub)
├── go.mod
├── Makefile                 # alvos run / build / test / tidy
└── README.md
```

Segue as convenções idiomáticas do ecossistema Go (`cmd/` para binários, `internal/` para lógica privada).

### Flags de `eco new`

- `--module <path>` — go module path (default: nome da pasta).
- `--force` — sobrescreve diretório não vazio.

## Desenvolvimento

```bash
git clone https://github.com/dixavier27/eco
cd eco
go build ./cmd/eco
go test ./...
```

## Licença

MIT. Veja [LICENSE](LICENSE).

---

<sub>O código pré-rewrite (TS/Bun, v2.x) está preservado na branch [`backup`](https://github.com/dixavier27/eco/tree/backup).</sub>
