# Jogo de Terminal em Go

Este projeto é um pequeno jogo desenvolvido em Go que roda no terminal usando a biblioteca [termbox-go](https://github.com/nsf/termbox-go). O jogador controla um personagem que pode se mover por um mapa carregado de um arquivo de texto.

## Como funciona

- O mapa é carregado de um arquivo `.txt` contendo caracteres que representam diferentes elementos do jogo.
- O personagem se move com as teclas **W**, **A**, **S**, **D**.
- Pressione **E** para interagir com o ambiente.
- Pressione **ESC** para sair do jogo.

### Controles

| Tecla | Ação              |
|-------|-------------------|
| W     | Mover para cima   |
| A     | Mover para esquerda |
| S     | Mover para baixo  |
| D     | Mover para direita |
| E     | Interagir         |
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


## Alterações feitas durante o trabalho

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