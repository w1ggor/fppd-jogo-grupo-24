## Relatório de Implementação — RPC, Cliente e Execução Única

### Visão geral

Este relatório descreve como o servidor RPC foi implementado, quais modificações foram feitas no cliente, e a estratégia adotada para garantir execução única (exactly-once) mesmo em cenários de retry e falhas temporárias de rede.

## Servidor RPC

O servidor está em `servidor/servidor.go` e expõe operações via `net/rpc` sobre TCP.

- Estado principal: a struct `DadosJogo` mantém o estado compartilhado do jogo (posições dos jogadores e flags de conexão) e é protegida por `sync.Mutex` para garantir segurança em concorrência.
- Cache de deduplicação: `comandosProcessados map[int]map[int]*ComandoProcessado` indexado por `ClientID` e `SequenceNumber`. Cada entrada guarda o resultado lógico do comando para viabilizar exactly-once.
- API RPC exposta:
	- `ConectarJogador(_ bool, numeroJogador *int)` — reserva slot e retorna `1`, `2` ou `-1` (cheio).
	- `Disconnect(clientId int, sucesso *bool)` — libera o slot do jogador.
	- `GetPosicoes(player int, resposta *Posicoes)` — devolve as posições correntes.
	- `Inicializar(dados Inicializar, sucesso *bool)` — define posições iniciais do jogo, com deduplicação por `SequenceNumber`.
	- `MoverJogador(dados MoverElementoType, sucesso *bool)` — atualiza a posição do jogador (1 ou 2), também com deduplicação.
- Concurrency: todo método RPC envolve `s.mu.Lock()/Unlock()` para serializar acesso a estado e ao cache de comandos.
- Loop de rede: `rpc.Register(servidor)` e `net.Listen` com `rpc.ServeConn` por conexão.

## Modificações no Cliente

O cliente principal (`main.go`) e a lógica de movimento (`personagem.go`) foram ajustados para usar RPC de forma resiliente e compatível com exactly-once:

- Retry com backoff: função `ChamarRPC(client, method, args, reply, maxRetries)` encapsula repetição com backoff simples em caso de erro temporário.
- Sequenciamento por cliente: em `jogo.go`, `Jogo` mantém `sequenceNumber` e o método `proximoSequenceNumber()` incrementa por comando lógico, garantindo unicidade local.
- Inicialização remota: após `jogoCarregarMapa`, o cliente envia `DadosJogo.Inicializar` contendo posições iniciais, `ClientID` (o número do jogador) e `SequenceNumber` gerado.
- Movimentação remota: em `personagemMover` (arquivo `personagem.go`), a cada movimento válido local, o cliente envia `DadosJogo.MoverJogador` com `{Player, X, Y, ClientID, SequenceNumber}`. Em caso de falha transitória, o retry é transparente.
- Sincronização contínua: uma goroutine no cliente chama periodicamente `DadosJogo.GetPosicoes` para refletir a posição do outro jogador e validar condições de vitória no cliente local.
- Desconexão graciosa: `Disconnect` é chamado em `defer` ao encerrar.

Estruturas RPC relevantes no cliente:

- `Posicoes { Pos1X, Pos1Y, Pos2X, Pos2Y, ClientID, SequenceNumber }`
- `MoverElementoTypeRPC { Player, X, Y, ClientID, SequenceNumber }`

Observação: existe também um `MoverElementoType` local (em `jogo.go`) usado apenas para atualização de tela/estado local; ele é distinto do tipo enviado via RPC.

## Estratégia de Execução Única (Exactly-Once)

Objetivo: permitir que o cliente faça retries (para robustez a perdas e timeouts) sem que o servidor aplique efeitos duplicados (ex.: mover o jogador duas vezes).

Componentes da estratégia:

1. Identidade do cliente: `ClientID` (1 ou 2) identifica o emissor.
2. Sequência por cliente: cada comando lógico gera um `SequenceNumber` crescente via `Jogo.proximoSequenceNumber()`.
3. Deduplicação no servidor:
	 - Antes de aplicar um comando, o servidor verifica em `comandosProcessados[ClientID][SequenceNumber]`.
	 - Se existir, retorna o `Resultado` cacheado e ignora reaplicação do efeito, logando mensagem de duplicado.
	 - Se não existir, aplica o efeito (ex.: atualiza `posicaoJogadores`), seta `*sucesso = true` e registra a entrada no cache.
4. Concorrência: acesso ao cache e ao estado é protegido por `Mutex` para evitar condições de corrida.

Propriedade alcançada:

- Com retry do cliente, o padrão base seria "pelo-menos-uma-vez". A tabela de deduplicação no servidor transforma o efeito em "exatamente-uma-vez" para cada par `(ClientID, SequenceNumber)`.

## Validação da Garantia Exactly-Once

Estratégia de validação usada no próprio código em execução:

- Geração de duplicatas: o cliente usa `ChamarRPC`, o que reemite chamadas em caso de erro; para fins de teste, basta provocar instabilidade (ou reduzir intervalos) para observar duplicatas.
- Evidência no servidor: os logs mostram linhas como:
	- "Movimentação do Jogador 1 para: (x, y) - SeqNum: N" quando o comando é aplicado.
	- "Movimento duplicado detectado - ClientID: 1, SeqNum: N (ignorado)" quando um retry chega com o mesmo `SequenceNumber`.
- Observável no cliente: a posição remota/local não avança duas vezes para o mesmo comando; o efeito visual é um único passo por tecla, sem "salto".
- Cobertura: verificado em `Inicializar` e `MoverJogador`. A obtenção de posições (`GetPosicoes`) é leitura e não exige deduplicação.