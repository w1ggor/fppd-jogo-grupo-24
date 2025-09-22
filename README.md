# Jogo de Terminal em Go

Este projeto é um pequeno jogo desenvolvido em Go que roda no terminal usando a biblioteca [termbox-go](https://github.com/nsf/termbox-go). O jogador controla um personagem que pode se mover por um mapa carregado de um arquivo de texto.

## Como funciona

- O mapa é carregado de um arquivo `.txt` contendo caracteres que representam diferentes elementos do jogo.
- O personagem de fogo se move com as teclas **W**, **A**, **S**, **D**.
- O personagem de água se move com as teclas **I**, **J**, **K**, **L**.
- Pressione **ESC** para sair do jogo.

### Controles Jogador 1

| Tecla | Ação              |
|-------|-------------------|
| W     | Mover para cima   |
| A     | Mover para esquerda |
| S     | Mover para baixo  |
| D     | Mover para direita |
| ESC   | Sair do jogo      |

### Controles Jogador 2

| Tecla | Ação              |
|-------|-------------------|
| I     | Mover para cima   |
| J     | Mover para esquerda |
| K     | Mover para baixo  |
| L     | Mover para direita |
| ESC   | Sair do jogo      |

## Como compilar

1. Instale o Go e clone este repositório.
2. Inicialize um novo módulo "jogo":

```bash
go mod init jogo
go get -u github.com/nsf/termbox-go
```

3. Compile o programa:

Linux:

```bash
go build -o jogo
```

Windows:

```bash
go build -o jogo.exe
```

Também é possivel compilar o projeto usando o comando `make` no Linux ou o script `build.bat` no Windows.

## Como executar

1. Certifique-se de ter o arquivo `mapa.txt` com um mapa válido.
2. Execute o programa no termimal:

```bash
./jogo
```

## Estrutura do projeto

- main.go — Ponto de entrada e loop principal
- interface.go — Entrada, saída e renderização com termbox
- jogo.go — Estruturas e lógica do estado do jogo
- personagem.go — Ações do jogador
- inimigo.go - ações dos inimigos


# Alterações feitas durante o trabalho
### Controle de dois jogadores
Foi adicionado um segundo jogador, dividindo eles em "fogo" e "água". Para controlar o input dos jogadores foi criada a função "recebeInput", que recebe os inputs via dois canais, um para cada jogador.
```go
// personagem.go
var player1Input = make(chan InputData)
var player2Input = make(chan InputData)

func recebeInput(player int, jogo *Jogo) {

	if player == 0 {
		for {
			var input = <-player1Input
			personagemMover(input, jogo, 0)
		}

	} else {
		for {

			var input = <-player2Input
			personagemMover(input, jogo, 1)
		}
	}
}

```

Na main, esta função é chamada uma vez para cada jogador, criando as duas goroutines que esperam pelos inputs via canais.

```go
// main.go
go recebeInput(0, &jogo)
go recebeInput(1, &jogo)
```
### Adição de inimigos com patrulha automática e canais concorrentes

Foram implementados dois inimigos no jogo: o inimigo de fogo (◇) e o inimigo de água (◆). Cada inimigo se move automaticamente pelo mapa, patrulhando de um lado para o outro.

- **Canais concorrentes:** Cada inimigo escuta dois canais: um canal de patrulha (para receber comandos de movimento) e um canal de alerta (para receber sinais de proximidade do personagem). Isso garante escuta concorrente de múltiplos canais.
- **Aceleração por proximidade:** Se o personagem de fogo se aproxima do inimigo de água, ou o personagem de água se aproxima do inimigo de fogo, o respectivo inimigo acelera sua movimentação automaticamente.
- **Colisão:** Se o inimigo de água encostar no personagem de fogo, o personagem de fogo retorna para sua posição inicial e a mensagem "Fogo apagou!" aparece na barra de status. Se o inimigo de fogo encostar no personagem de água, o personagem de água retorna para sua posição inicial e a mensagem "Agua ferveu!" aparece na barra de status.
- **Sincronização:** Toda movimentação dos inimigos também é feita via canal, garantindo concorrência segura.

Essas alterações demonstram comunicação entre elementos do jogo por canais, escuta concorrente e lógica reativa baseada em eventos do ambiente.
### Botões que abrem e fecham portões 
Foram implementados dois botões e dois portões interativos em cada um dos lados do mapa. Quando o jogador fica em cima do botão de um lado, o portão do lado oposto irá abrir. Quando ele sai, o portão irá fechar novamente.

O portão fechará apenas após abrir por completo.
- **Concorrência:** Na main, é chamado um método que inicia ambos os botões como goroutines que ficam esperando até que um jogador fique em cima de um deles. Quando precionado, irá ativar, de maneira concorrente, outra goroutine, que é encarregada por abrir o portão, e logo em seguida uma terceira, com a função de fechar o portão.
- **Canais:** Também são utilizados quatro canais, sendo cada um para controlar a abertura e fechamento de cada portão. O método principal de cada botão ficará esperando até que a goroutine de abrir o portão envie uma mensagem pelo canal indicando que o portão abriu. Somente então ele poderá iniciar o fechamento. Da mesma forma, a próxima abertura espera o portão fechar por completo.
### Bandeiras que finalizam o jogo
Foram implementadas duas bandeiras, uma para cada jogador, para marcar a condição de vitória do jogo, sendo assim, ambos os jogadores precisam estar nas bandeiras ao mesmo tempo para ganharem.
- **Canais concorrentes**: Junto à chamada da função das bandeiras, é iniciada uma goroutine para mostrar o aviso de 15 segundos faltando para acabar.
- **Escuta múltiplos canais**: Em jogoPodeMoverPara, são enviadas mensagens nos canais player1Vence e player2Vence, para indicar quando cada jogador está em cima de uma bandeira. Em vencerJogo, no código das bandeiras, há um select que espera pela mensagem destes dois canais, para verificar se ambos chegaram na bandeira e avisar que ganharam o jogo.
- **Timeout**: Junto ao select de vencerJogo, há um timeout de 30 segundos para avisar os jogadores que perderam e reiniciar o jogo. Este timeout é implementado com time.After(), em formato que canal, que é "recebido" pelo select.

### Modificamos o sistema de atualização da interface
Modificamos a forma que o jogo atualiza o mapa. Antes, a função "interfaceDesenharJogo" era chamada em função da movimentação do jogador, fizemos com que ela fosse chamada aproximadamente 60 vezes por segundo dentro de uma goroutine, feita por meio de uma função anônima.

```go
//main.go
go func() {
	for {
		// Verifica colisão inimigo de água com personagem de fogo
		if jogo.IniAguaPosX == jogo.Pos1X && jogo.IniAguaPosY == jogo.Pos1Y {
			// Volta personagem de fogo para posição inicial
			apagarFogo(&jogo)
		}
		// Verifica colisão inimigo de fogo com personagem de água
		if jogo.IniFogoPosX == jogo.Pos2X && jogo.IniFogoPosY == jogo.Pos2Y {
			// Volta personagem de água para posição inicial
			evaporarAgua(&jogo)
		}
		interfaceDesenharJogo(&jogo)
		time.Sleep(16 * time.Millisecond)
	}
}()
```
### Agua e lava
Os dois jogadores ficaram divididos em Água e Lava. Foram adicionados elemetos no mapa que interagem apenas com um dos jogadores, as barreiras de água impedem o jogador "fogo" de passar e as barreiras de fogo
### Foi adicionado um canal para sincronizar a mudança do mapa
Fizemos a sincronização da atualização do mapa via um canal com buffer de tamanho 1, garantindo assim que apenas um elemento pode atualizar o mapa por vez.

<details>
<summary>Adicionamos uma struct com as informações necessárias para atualizar o mapa</summary>

```go
//jogo.go
type MoverElementoType struct {
	jogo         *Jogo
	x, y, dx, dy int
}
```
</details>

<details>
<summary>Depois criamos um canal que vai receber estas informações.</summary>

```go
//jogo.go
var moveElemento = make(chan MoverElementoType, 1)
```
</details>

<details>
<summary>Então removemos os parâmetros da função "jogoMoverElemento", agora recebendo estas informações através do canal "moveElemento".</summary>

```go
// jogo.go
func jogoMoverElemento() {
	for {
		var moveInput = <-moveElemento
		var jogo = moveInput.jogo
		var x, y, dx, dy = moveInput.x, moveInput.y, moveInput.dx, moveInput.dy
		nx, ny := x+dx, y+dy
	
		elemento := jogo.Mapa[y][x] 
		jogo.Mapa[y][x] = jogo.UltimoVisitado  
		jogo.UltimoVisitado = jogo.Mapa[ny][nx] 
		jogo.Mapa[ny][nx] = elemento
	}
}
```
</details>

# Requisitos do trabalho

## Foram implementados ao menos 3 tipos de elementos concorrentes autônomos com comportamentos visíveis e distintos no mapa
- Dois jogadores
- Botões e portões
- Inimigos
- Bandeiras
- Água e Lava

## Há uso de canais para comunicação e sincronização entre elementos 
- Sistema de inputs
- Sincronização do mapa
## Pelo menos um elemento escuta múltiplos canais e isso é demonstrável
- Verificação das bandeiras para vencer o jogo
## Pelo menos um elemento utiliza canais com timeout de forma testável
- Timeout que finaliza o jogo e mostra a mensagem
## Há controle de exclusão mútua nas regiões críticas do jogo utilizando canais 
- Sincronização do mapa
- Controle dos portões