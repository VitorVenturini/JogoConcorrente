# Arena Concorrente 🛡️

Projeto desenvolvido por Rodrigo e Vitor em linguagem Go para a disciplina de **Fundamentos de Processamento Paralelo e Distribuido (FPPD)**.

A aplicacao consiste em um jogo interativo de terminal projetado inteiramente sob uma arquitetura orientada a goroutines e channels, aplicando de forma pratica a teoria de concorrencia com passagem de mensagens, eliminando o uso de travas na memoria.

## Caracteristicas do Projeto Finalizado

- **Arena Expansiva**: Grid de combate dimensionado em 25x25.
- **Alta Concorrencia**: 10 goroutines rodando simultaneamente sem `sync.Mutex`.
- **I.A. Independente**: 6 inimigos autonomos processados de forma assincrona com base em tempo (`Ticker`) individual, garantindo total independencia das acoes do jogador.
- **I.A. com Caminho**: Inimigos usam busca em largura (BFS) para contornar obstaculos e perseguir o jogador.
- **Obstaculos**: Paredes aleatorias no grid, bloqueando movimento do jogador e dos inimigos.
- **Spawn Seguro**: Jogador e inimigos nunca nascem sobre obstaculos.
- **Ataque de Alvo Unico**: Cada ataque afeta apenas um inimigo adjacente por vez.
- **Velocidade Ajustada**: Inimigos com intervalos menores para aumentar a pressao do jogo.
- **Performance Grafica**: Renderizacao implementada com a biblioteca `tcell` para garantir leitura e plotagem em "Raw Mode", mitigando travamentos e o classico *flickering* (cintilacao) de terminais comuns.
- **Shutdown Limpo**: Todo o ciclo de vida das goroutines e dos recursos de terminal sao fechados graciosamente atraves do `context.WithCancel` e sincronia do `sync.WaitGroup`.

---

## Como Configurar e Executar

Este projeto depende da biblioteca grafica de terminal *tcell*. Antes de rodar, e necessario instalar as dependencias.

Abra o seu terminal (preferencialmente PowerShell ou terminal nativo do SO) e execute:

```powershell
# 1. Instalar a biblioteca de interface de terminal
go get github.com/gdamore/tcell/v2

# 2. Executar o jogo
go run .
```

## Como Validar a Arquitetura (Race Conditions)

O projeto foi construido do zero sob as premissas de boas praticas do Go, canalizando a propriedade mutavel em vez de compartilhar ponteiros inseguros. Para provar a robustez em tempo de avaliacao, execute o jogo com a flag do detector de corridas do compilador Go:

```powershell
go run -race .
```

## Diagnostico de Travamentos e Panics

Go nao possui try/catch. Para diagnosticar travamentos e panics, o jogo usa `defer` + `recover` nas goroutines e grava logs em arquivo.

Sempre que o jogo rodar, um arquivo `debug.log` sera criado/atualizado na pasta do projeto. Caso ocorra um travamento, feche o jogo e verifique as ultimas linhas desse arquivo para identificar o componente e a causa.

## Controles do Jogo

- ⬅️ ⬆️ ⬇️ ➡️ **Setas Direcionais**: Movem o jogador (🙂) livremente pelo grid.
- ⎵ **Barra de Espaco**: Executa um ataque fisico direto (atinge apenas um inimigo adjacente por vez).
- ⎋ **Tecla Esc**: Aciona a CancelFunc global, enviando sinal de shutdown gracioso e fechando o jogo a qualquer momento (inclusive na tela de Game Over/Vitoria).

## Notas de Renderizacao

Os personagens usam emojis. Para evitar sobreposicao visual com os obstaculos, a arena e renderizada com largura dupla por celula, mantendo cada emoji dentro da sua casa.

## Destaques Academicos da Entrega

- O uso do bloco de controle `select` como ferramenta multiplexadora atomica de entrada de dados, lidando com 5 ramificacoes de eventos simultaneos.
- Transformacao do sistema em tempo-real (semelhante a um main loop convencional) sem perder as abstracoes de design de sistemas paralelos.