# Arena Concorrente

Projeto em Go para a disciplina de programação concorrente. A proposta é um jogo/simulação de terminal com arquitetura orientada a goroutines e channels.

## Estado atual

Base inicial do projeto criada:

- `go.mod`
- `main.go`
- `types.go`
- `coordinator.go`
- `renderer.go`
- `input.go`
- documentação em `ARQUITETURA.md`
- backlog em `BACKLOG.md`

## Como executar

```powershell
go run .
```

## Como validar race conditions

```powershell
go run -race .
```

## Próximas etapas

- ticker
- inimigos autônomos
- shutdown gracioso

## Estado atual do fluxo

O projeto já possui:

- goroutine de input;
- goroutine de renderização;
- goroutine coordenadora com `select`;
- movimentação básica do jogador;
- comando `quit` para encerrar o loop principal.
