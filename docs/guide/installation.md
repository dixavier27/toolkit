# Instalação

## Via npm/bun (recomendado)

::: code-group

```bash [Bun]
bun add -D github:dixavier27/toolkit#v2.8.0
```

```bash [pnpm]
pnpm add -D github:dixavier27/toolkit#v2.8.0
```

```bash [npm]
npm install -D github:dixavier27/toolkit#v2.8.0
```

:::

O pacote é publicado como `@dixavier27/eco` (o repositório se chama `toolkit` por motivos históricos).

## Global (para usar `eco` direto no terminal)

```bash
bun add -g github:dixavier27/toolkit#v2.8.0
```

Depois, em qualquer projeto:

```bash
eco --version
eco --help
```

## Via bunx (sem instalar)

Para uso pontual (ex: scaffolar novo projeto):

```bash
bunx @dixavier27/eco new minha-app
```

## Verificação

```bash
eco info
```

Saída esperada:

```
eco — informações do ambiente

  eco                  2.8.0
  Bun                  1.3.x
  javascript-obfuscator 5.4.x
  Plataforma host      linux (x64)
  Config detectado     nenhum (rode 'eco init')
```

Se algo estiver faltando, rode `eco doctor --fix`.
