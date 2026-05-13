# Backlog do Projeto

## Objetivo

Este backlog dividiu o projeto em etapas pequenas, incrementais e testaveis. A ideia central foi permitir que o grupo implementasse o jogo de forma progressiva, validando a concorrencia, comunicacao por channels, renderizacao e shutdown gracioso em cada fase, ate culminar no estado atual de entrega.

Projeto desenvolvido por Rodrigo e Vitor.

O projeto de referencia deste backlog e o detalhado em `ARQUITETURA.md`: **Arena Concorrente**.

---

## Estrategia de Execucao

### Principios Adotados

- implementar primeiro a espinha dorsal concorrente;
- manter o estado restrito e centralizado no `game coordinator`;
- adicionar uma capacidade nova por etapa;
- garantir validacao constante atraves de testes com a flag `-race`.

### Definicao de "Pronto" (Definition of Done)

O projeto e considerado pronto porque:

- [x] compila com `go run .`;
- [x] executa sem erros e sem alertas com `go run -race .`;
- [x] possui pelo menos 4 goroutines com papeis distintos;
- [x] possui pelo menos 2 (foram entregues 6) entidades autonomas com goroutines proprias;
- [x] usa channels como mecanismo exclusivo de comunicacao do estado;
- [x] usa dezenas de ramificacoes `select` no coordenador e agentes;
- [x] encerra sem *goroutine leak* via contexto;
- [x] renderiza em tempo real de forma polida (via `tcell`);
- [x] possui documentacao arquitetural coerente com o codigo final.

---

## Visao Geral das Etapas Concluidas

1. Inicializacao do projeto Go e dominio.
2. Coordenador central com loop e `select` estatico.
3. Migracao da renderizacao e input para biblioteca `tcell` (raw mode).
4. Ticker temporal do jogo.
5. Logica de multiplos inimigos autonomos (A, B, C e D).
6. Regras de colisao, dano e travamento de tela de fim de jogo.
7. Polimento, shutdown gracioso e documentacao final.

---

## Historico do Backlog de Implementacao (Concluido)

### Etapa 1 e 2. Base e Dominio Go

- **Tarefas:** Criacao do `go.mod`, `main.go`, `types.go`.
- **Resultado:** Definicao dos structs `PlayerCommand`, `EnemyAction`, `Tick`, `GameSnapshot` e `GameState`. Tipagem rigorosa para trafegar dados com seguranca nos channels.

### Etapa 3. Coordenador Central e select

- **Tarefas:** Criacao da goroutine central `coordinator.go` e dos channels base.
- **Resultado:** A goroutine assumiu seu papel de receber input, processar e devolver snapshots. Eliminacao da necessidade do uso de travas na memoria de entidades.

### Etapa 4 e 5. Renderer e Input via tcell

- **Tarefas:** Descarte da abordagem bloqueante (`fmt.Scan`) para o uso da biblioteca `github.com/gdamore/tcell/v2`.
- **Resultado:** O jogo ganhou *double-buffering* e captura instantanea das teclas Setas (Movimento), Espaco (Ataque) e Esc (Quit), mudando a sensacao de "simulacao turno-a-turno" para "tempo real responsivo".

### Etapa 6. Tempo Concorrente

- **Tarefas:** Goroutine de pulsos `ticker.go`.
- **Resultado:** Insercao de uma batida constante que alimenta o Coordenador, servindo de base para a HUD sem acoplar a velocidade do ambiente ao teclado do jogador.

### Etapa 7 e 8. Entidades Autonomas (Inimigos)

- **Tarefas:** Criacao do script `enemy.go` generico e instanciamento de 6 goroutines independentes na `main`.
- **Resultado:** Inimigos com algoritmos de Chase (Perseguicao) e Patrol (Patrulha). Cada um ganhou um tempo de processamento distinto (`speed time.Duration`), agindo sob seus proprios ritmos independentes do jogador, cumprindo estritamente a exigencia do professor.

### Etapa 9. Colisao, Dano e State Locking

- **Tarefas:** Logica de intersecao de coordenadas `X` e `Y`.
- **Resultado:** Jogador causa dano com `Espaco` e inimigos tiram vida ao pisar na mesma casa. Adicionado um bloqueio no Coordenador: se `GameOver` ou `Victory` ocorrerem, os canais ignoram movimentacoes futuras e esperam estaticamente a acao de finalizacao (Esc).

### Etapa 10 e 11. Shutdown Gracioso e Expansao

- **Tarefas:** Blindar a aplicacao contra vazamentos e escalar os desafios.
- **Resultado:** A arena cresceu para 25x25, com HUD informativo no rodape. Todos os blocos `for` escutam `ctx.Done()` e param seus *tickers*. A `main` so chama a finalizacao do terminal apos `wg.Wait()` atestar o fim seguro de todos os 10 processos concorrentes.

### Etapa 12. Documentacao e Preparacao para Bancas

- **Tarefas:** Atualizar `ARQUITETURA.md` e `README.md`.
- **Resultado:** Artefatos prontos e codigo imune no teste `go run -race .`. Preparacao conceitual de respostas para questionamentos dos professores em cima do bloco `select` e da exclusao intencional de `Mutex`.

---

## Checklist Executivo (Finalizado)

- [x] Definir arquitetura concorrente centralizada em eventos.
- [x] Construir Coordenador, tipos e canais.
- [x] Implementar e refinar *tcell* no Input/Render.
- [x] Lancar 6 entidades IA autonomas assincronas.
- [x] Adicionar regras e travamento de Game Over.
- [x] Assegurar shutdown gracioso com `context`.
- [x] Validar 100% livre de falhas no Race Detector.
- [x] Finalizar documentacao.
