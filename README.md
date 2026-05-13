# Arena Concorrente 🛡️

Projeto desenvolvido por Rodrigo e Vitor em linguagem Go para a disciplina de **Fundamentos de Processamento Paralelo e Distribuido (FPPD)**.

A aplicacao consiste em um jogo interativo de terminal projetado inteiramente sob uma arquitetura orientada a goroutines e channels, aplicando de forma pratica a teoria de concorrencia com passagem de mensagens, eliminando o uso de travas na memoria.

## Caracteristicas do Projeto Finalizado

- **Arena Expansiva**: Grid de combate dimensionado em 25x25.
- **Alta Concorrencia**: 10 goroutines rodando simultaneamente sem `sync.Mutex`.
- **I.A. Independente**: 6 inimigos autonomos processados de forma assincrona com base em tempo (`Ticker`) individual, garantindo total independencia das acoes do jogador.
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

## Controles do Jogo

- ⬅️ ⬆️ ⬇️ ➡️ **Setas Direcionais**: Movem o jogador (@) livremente pelo grid.
- ⎵ **Barra de Espaco**: Executa um ataque fisico direto (atinge simultaneamente todos os inimigos localizados nas casas em cruz adjacentes).
- ⎋ **Tecla Esc**: Aciona a CancelFunc global, enviando sinal de shutdown gracioso e fechando o jogo a qualquer momento (inclusive na tela de Game Over/Vitoria).

## Destaques Academicos da Entrega

- O uso do bloco de controle `select` como ferramenta multiplexadora atomica de entrada de dados, lidando com 5 ramificacoes de eventos simultaneos.
- Transformacao do sistema em tempo-real (semelhante a um main loop convencional) sem perder as abstracoes de design de sistemas paralelos.