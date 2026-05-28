# eco

CLI em Go para gerenciar ambientes de desenvolvimento de APIs REST em Go.

> Reinício enxuto. O histórico anterior (TS/Bun, v2.x) está preservado na branch [`backup`](https://github.com/dixavier27/eco/tree/backup).

## Instalação

```bash
go install github.com/dixavier27/eco/cmd/eco@latest
```

## Uso

```bash
eco version          # imprime a versão
eco doctor           # verifica toolchain Go
eco new meu-app      # cria nova API REST em ./meu-app
```

## Status

MVP em construção. Comandos planejados (`build`, `run`, `lint`, `test`, `release`, `ci`) serão reintroduzidos gradualmente — ver plano de desenvolvimento.
