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
### Adição de inimigos com patrulha automática e canais concorrentes

Foram implementados dois inimigos no jogo: o inimigo de fogo (◇) e o inimigo de água (◆). Cada inimigo se move automaticamente pelo mapa, patrulhando de um lado para o outro.

- **Canais concorrentes:** Cada inimigo escuta dois canais: um canal de patrulha (para receber comandos de movimento) e um canal de alerta (para receber sinais de proximidade do personagem). Isso garante escuta concorrente de múltiplos canais.
- **Aceleração por proximidade:** Se o personagem de fogo se aproxima do inimigo de água, ou o personagem de água se aproxima do inimigo de fogo, o respectivo inimigo acelera sua movimentação automaticamente.
- **Colisão:** Se o inimigo de água encostar no personagem de fogo, o personagem de fogo retorna para sua posição inicial e a mensagem "Fogo apagou!" aparece na barra de status. Se o inimigo de fogo encostar no personagem de água, o personagem de água retorna para sua posição inicial e a mensagem "Agua ferveu!" aparece na barra de status.
- **Sincronização:** Toda movimentação dos inimigos também é feita via canal, garantindo concorrência segura.

Essas alterações demonstram comunicação entre elementos do jogo por canais, escuta concorrente e lógica reativa baseada em eventos do ambiente.
### Botões que abrem e fecham portões 
[TODO]
### Bandeiras que finalizam o jogo
[TODO]
### Modificamos o sistema de atualização da interface
### Agua e lava

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