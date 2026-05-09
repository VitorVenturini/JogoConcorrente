# Backlog do Projeto

## Objetivo

Este backlog divide o projeto em etapas pequenas, incrementais e testáveis. A ideia é permitir que o grupo implemente o jogo de forma progressiva, validando concorrência, comunicação por channels, renderização e shutdown gracioso em cada fase.

O projeto de referência deste backlog é o proposto em [ARQUITETURA.md](/c:/repositorios/JogoConcorrente/ARQUITETURA.md): **Arena Concorrente**.

---

## Estratégia de Execução

### Princípios

- implementar primeiro a espinha dorsal concorrente;
- manter o estado centralizado no `game coordinator`;
- adicionar uma capacidade nova por etapa;
- garantir que toda etapa tenha um teste manual simples;
- rodar `go run -race` assim que o código começar a existir.

### Definição de pronto do projeto

O projeto estará pronto quando:

- compilar com `go run .`;
- executar sem erros com `go run -race .`;
- possuir pelo menos 4 goroutines com papéis distintos;
- possuir pelo menos 2 entidades autônomas com goroutines próprias;
- usar channels como mecanismo principal de comunicação;
- usar `select` no coordenador;
- encerrar sem goroutine leak;
- renderizar em tempo real no terminal;
- possuir documento de arquitetura coerente com o código.

---

## Visão Geral das Etapas

1. Inicialização do projeto Go
2. Modelo de domínio e tipos de evento
3. Coordenador central com loop e `select`
4. Renderização básica no terminal
5. Input por comando digitado
6. Ticker do jogo
7. Primeiro inimigo autônomo
8. Segundo inimigo autônomo
9. Regras de colisão, dano e fim de jogo
10. Shutdown gracioso e robustez
11. Polimento de jogabilidade
12. Documentação e preparação para apresentação

---

## Backlog por Etapa

## Etapa 1. Inicialização do projeto Go

### Objetivo
Criar a base mínima do projeto para permitir evolução incremental.

### Tarefas

- criar `go.mod`;
- criar `main.go`;
- definir estrutura inicial de pastas;
- decidir se o projeto ficará todo em `package main` no início ou separado em pacotes;
- adicionar um `README.md` curto com instruções de execução.

### Entregáveis

- `go.mod`
- `main.go`
- estrutura inicial de diretórios

### Critérios de aceite

- `go run .` executa sem erro;
- o programa imprime mensagem inicial ou abre loop mínimo do jogo.

### Como testar

```powershell
go run .
```

### Observações

Para este trabalho, começar com estrutura simples é melhor do que tentar modularizar cedo demais.

---

## Etapa 2. Modelo de domínio e tipos de evento

### Objetivo
Definir os tipos de dados que circularão entre as goroutines.

### Tarefas

- definir `PlayerCommand`;
- definir `EnemyAction`;
- definir `Tick`;
- definir `GameSnapshot`;
- definir `GameState`;
- definir tipos de posição, arena, entidades e HUD.

### Entregáveis

- arquivo com tipos de domínio, por exemplo `types.go`

### Critérios de aceite

- os tipos permitem representar jogador, inimigos, mapa, eventos e renderização;
- o estado do jogo é suficientemente claro para o coordenador controlar tudo sozinho.

### Como testar

- compilar o projeto e confirmar que os tipos estão coerentes;
- opcionalmente criar um `main` temporário que instancie um `GameState` de exemplo.

```powershell
go run .
```

---

## Etapa 3. Coordenador central com loop e `select`

### Objetivo
Implementar a peça principal da arquitetura concorrente.

### Tarefas

- criar a goroutine `game coordinator`;
- criar os channels principais;
- implementar loop com `select`;
- tratar casos de input, ação de inimigo, tick e cancelamento;
- gerar snapshots iniciais para o renderer.

### Entregáveis

- arquivo como `coordinator.go`
- canais criados e conectados na `main`

### Critérios de aceite

- o coordenador roda em goroutine própria;
- existe `select` multiplexando pelo menos `inputCh`, `enemyActionCh`, `tickCh` e `ctx.Done()`;
- não há acesso concorrente desnecessário ao estado principal.

### Como testar

- subir o programa com eventos simulados enviados pela `main`;
- verificar por logs ou prints se o coordenador processa cada evento corretamente.

```powershell
go run .
go run -race .
```

### Dependências

- etapa 1
- etapa 2

---

## Etapa 4. Renderização básica no terminal

### Objetivo
Mostrar o estado do jogo em tempo real no terminal.

### Tarefas

- criar goroutine `renderer`;
- implementar limpeza de tela com ANSI;
- desenhar grade simples da arena;
- desenhar jogador e inimigos a partir de `GameSnapshot`;
- renderizar HUD mínima com vida, turno ou tempo.

### Entregáveis

- arquivo como `renderer.go`

### Critérios de aceite

- apenas o renderer escreve na saída visual do jogo;
- a tela é redesenhada a partir de snapshots;
- o jogo pode mostrar um frame estático sem concorrência incorreta.

### Como testar

- iniciar com estado fixo;
- enviar snapshots artificiais pelo `renderCh`;
- validar se a arena aparece corretamente no terminal.

```powershell
go run .
```

### Dependências

- etapa 2
- etapa 3

---

## Etapa 5. Input por comando digitado

### Objetivo
Permitir controlar o jogador com entrada simples e apresentável.

### Tarefas

- criar goroutine `input`;
- ler comandos por linha;
- mapear texto para `PlayerCommand`;
- aceitar comandos como `up`, `down`, `left`, `right`, `attack`, `quit`;
- enviar comandos ao coordenador por `inputCh`.

### Entregáveis

- arquivo como `input.go`

### Critérios de aceite

- o usuário consegue mover o jogador;
- comando `quit` encerra a aplicação;
- comandos inválidos não derrubam o programa.

### Como testar

```powershell
go run .
```

Testes manuais:

- digitar `up`;
- digitar `left`;
- digitar `attack`;
- digitar `quit`.

### Dependências

- etapa 2
- etapa 3
- etapa 4

---

## Etapa 6. Ticker do jogo

### Objetivo
Introduzir passagem de tempo concorrente.

### Tarefas

- criar goroutine `ticker`;
- enviar eventos periódicos para `tickCh`;
- definir frequência de atualização;
- fazer o coordenador reagir a cada tick.

### Entregáveis

- arquivo como `ticker.go`

### Critérios de aceite

- ticks chegam regularmente ao coordenador;
- o estado do jogo evolui sem depender apenas de input manual.

### Como testar

- exibir contador de ticks no HUD;
- verificar se o contador aumenta automaticamente.

```powershell
go run .
go run -race .
```

### Dependências

- etapa 3

---

## Etapa 7. Primeiro inimigo autônomo

### Objetivo
Adicionar a primeira entidade concorrente independente.

### Tarefas

- criar goroutine do `enemy A`;
- definir lógica simples de movimento, por exemplo perseguição ou movimento aleatório;
- enviar intenções ao coordenador via `enemyActionCh`;
- fazer o coordenador aplicar a ação no estado.

### Entregáveis

- arquivo como `enemy.go` ou `enemy_a.go`

### Critérios de aceite

- o inimigo se move sem input do jogador;
- a lógica roda em goroutine própria;
- o inimigo não altera estado global diretamente.

### Como testar

- iniciar o jogo sem mover o jogador;
- observar o inimigo se movendo sozinho;
- rodar com detector de race.

```powershell
go run .
go run -race .
```

### Dependências

- etapa 3
- etapa 4
- etapa 6

---

## Etapa 8. Segundo inimigo autônomo

### Objetivo
Cumprir explicitamente o requisito de duas entidades autônomas.

### Tarefas

- criar goroutine do `enemy B`;
- dar comportamento próprio, por exemplo patrulha, perseguição mais lenta ou ataque à distância;
- integrar ao mesmo fluxo de eventos do coordenador.

### Entregáveis

- segundo inimigo funcionando com goroutine própria

### Critérios de aceite

- existem duas entidades autônomas independentes;
- ambas enviam ações por channels;
- seus comportamentos são observáveis e distintos.

### Como testar

- iniciar o jogo e observar os dois inimigos agindo ao mesmo tempo;
- validar se continuam funcionando durante movimentação do jogador.

```powershell
go run .
go run -race .
```

### Dependências

- etapa 7

---

## Etapa 9. Regras de colisão, dano e fim de jogo

### Objetivo
Transformar a simulação em jogo completo.

### Tarefas

- implementar colisão com bordas da arena;
- implementar colisão entre entidades;
- implementar dano, vida e derrota;
- implementar condição de vitória, por exemplo sobreviver por N ticks ou derrotar inimigos;
- exibir mensagens de fim de jogo.

### Entregáveis

- regras de jogo completas no coordenador

### Critérios de aceite

- o jogo tem começo, meio e fim;
- as regras estão centralizadas e previsíveis;
- o coordenador decide vitória e derrota.

### Como testar

- provocar colisões intencionais;
- deixar os inimigos atingirem o jogador;
- validar vitória e derrota em cenários simples.

```powershell
go run .
```

### Dependências

- etapa 5
- etapa 6
- etapa 7
- etapa 8

---

## Etapa 10. Shutdown gracioso e robustez

### Objetivo
Garantir encerramento limpo e cumprir um dos requisitos mais importantes da avaliação.

### Tarefas

- usar `context.WithCancel` ou `doneCh`;
- integrar `sync.WaitGroup`;
- fazer todas as goroutines observarem cancelamento;
- encerrar timers e loops corretamente;
- evitar envios em channels após shutdown.

### Entregáveis

- ciclo de encerramento completo

### Critérios de aceite

- `quit` encerra o jogo sem travar;
- fim de vitória ou derrota encerra o jogo sem vazamento;
- `go run -race .` não acusa race;
- não há goroutines presas visivelmente.

### Como testar

```powershell
go run .
go run -race .
```

Testes manuais:

- iniciar e sair com `quit`;
- perder o jogo e confirmar encerramento;
- vencer o jogo e confirmar encerramento.

### Dependências

- etapas 3 a 9

---

## Etapa 11. Polimento de jogabilidade

### Objetivo
Melhorar a experiência sem comprometer a simplicidade da apresentação.

### Tarefas

- ajustar tamanho da arena;
- balancear velocidade dos inimigos;
- melhorar HUD;
- adicionar mensagens de ajuda na tela;
- refinar símbolos visuais de jogador, inimigos e obstáculos;
- reduzir flicker de renderização, se necessário.

### Entregáveis

- jogo mais legível e demonstrável

### Critérios de aceite

- o jogo está fácil de entender visualmente;
- a demonstração ao vivo fica clara;
- os elementos importantes aparecem no terminal sem confusão.

### Como testar

- executar sessões curtas de jogo;
- verificar se a apresentação visual ajuda e não atrapalha.

```powershell
go run .
```

### Dependências

- etapas 4 a 10

---

## Etapa 12. Documentação e preparação para apresentação

### Objetivo
Fechar a entrega acadêmica e alinhar o código ao documento.

### Tarefas

- revisar [ARQUITETURA.md](/c:/repositorios/JogoConcorrente/ARQUITETURA.md) conforme a implementação final;
- converter o conteúdo para PDF depois;
- documentar comando de execução e teste com `-race`;
- preparar roteiro curto da apresentação;
- revisar nomes de goroutines, channels e fluxo real do código.

### Entregáveis

- documento final coerente com a implementação
- repositório pronto para entrega

### Critérios de aceite

- o diagrama bate com o código;
- o grupo consegue explicar cada goroutine e channel;
- o professor consegue rodar o projeto.

### Como testar

- fazer ensaio de apresentação;
- pedir para alguém do grupo explicar a arquitetura olhando só para o diagrama;
- validar execução do zero em outra máquina, se possível.

### Dependências

- todas as etapas anteriores

---

## Backlog Técnico Transversal

Esses itens devem ser acompanhados durante várias etapas, não apenas no final.

### Concorrência

- garantir que o estado principal seja modificado só pelo coordenador;
- evitar mutex, exceto se houver necessidade real de renderização;
- preferir snapshots imutáveis para a saída visual.

### Qualidade

- rodar `go run -race .` com frequência;
- manter funções curtas;
- evitar lógica de regra de jogo espalhada em várias goroutines.

### Observabilidade

- enquanto o projeto estiver em construção, usar logs simples para depurar eventos;
- remover ou reduzir logs antes da versão final de apresentação.

### Apresentação

- manter nomes simples e didáticos;
- evitar features extras que dificultem explicar a arquitetura;
- privilegiar clareza em vez de sofisticação desnecessária.

---

## Sequência Recomendada de Implementação

Se o grupo quiser seguir a ordem mais segura, use esta:

1. Etapa 1
2. Etapa 2
3. Etapa 3
4. Etapa 4
5. Etapa 6
6. Etapa 5
7. Etapa 7
8. Etapa 8
9. Etapa 9
10. Etapa 10
11. Etapa 11
12. Etapa 12

Essa ordem funciona bem porque primeiro cria o núcleo concorrente, depois a visualização, depois o tempo, depois o input, e por fim as entidades autônomas e regras.

---

## Checklist Executivo

### Planejamento

- [x] Definir proposta do jogo
- [x] Definir arquitetura concorrente
- [x] Criar documento de arquitetura
- [x] Criar backlog incremental

### Implementação

- [x] Inicializar módulo Go
- [x] Criar tipos do domínio
- [x] Criar coordinator com `select`
- [x] Criar renderer
- [x] Criar input
- [ ] Criar ticker
- [ ] Criar enemy A
- [ ] Criar enemy B
- [ ] Implementar colisões e dano
- [ ] Implementar vitória e derrota
- [ ] Implementar shutdown gracioso
- [ ] Rodar sem race

### Entrega

- [ ] Revisar documentação
- [ ] Preparar PDF
- [ ] Revisar repositório
- [ ] Ensaiar apresentação

---

## Próximo Passo Recomendado

O próximo passo de implementação deveria ser:

1. criar `go.mod`;
2. criar `main.go`;
3. definir tipos do domínio;
4. subir o coordenador com `select` e um renderer mínimo.

Esse corte já cria a base técnica mais importante do projeto e permite testar o fluxo principal antes de adicionar inimigos e regras de jogo.
