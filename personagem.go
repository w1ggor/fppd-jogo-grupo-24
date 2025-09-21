// personagem.go - Funções para movimentação e ações do personagem
package main

import (
	"fmt"
)

// Atualiza a posição do personagem com base na tecla pressionada (WASD)
func personagemMover(input InputData, jogo *Jogo, player int) {

	dx, dy := input.dx, input.dy

	if player == 0 {
		nx, ny := jogo.Pos1X+dx, jogo.Pos1Y+dy
		// Verifica se o movimento é permitido e realiza a movimentação
		if jogoPodeMoverPara(jogo, nx, ny) {
			var moveInput = MoverElementoType{player: 0, jogo: jogo, x: jogo.Pos1X, y: jogo.Pos1Y, dx: dx, dy: dy}
			moveElemento <- moveInput
			jogo.Pos1X, jogo.Pos1Y = nx, ny
		}
	} else {
		nx, ny := jogo.Pos2X+dx, jogo.Pos2Y+dy
		// Verifica se o movimento é permitido e realiza a movimentação
		if jogoPodeMoverPara(jogo, nx, ny) {
			var moveInput = MoverElementoType{player: 1,jogo: jogo, x: jogo.Pos2X, y: jogo.Pos2Y, dx: dx, dy: dy}
			moveElemento <- moveInput
			jogo.Pos2X, jogo.Pos2Y = nx, ny
		}
	}

}

// Define o que ocorre quando o jogador pressiona a tecla de interação
// Neste exemplo, apenas exibe uma mensagem de status
// Você pode expandir essa função para incluir lógica de interação com objetos
func personagemInteragir(jogo *Jogo) {
	// Atualmente apenas exibe uma mensagem de status
	jogo.StatusMsg = fmt.Sprintf("Interagindo em (%d, %d)", jogo.Pos1X, jogo.Pos1Y)
}

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

// Processa o evento do teclado e executa a ação correspondente
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo) bool {
	var input = InputData{player: 0, input: ev, dx: 0, dy: 0}
	switch ev.Tipo {
	case "sair":
		// Retorna false para indicar que o jogo deve terminar
		return false
	case "interagir":
		// Executa a ação de interação
		personagemInteragir(jogo)
	case "mover":
		// Move o personagem com base na tecla
		switch ev.Tecla {
		case 'w':
			input.player = 0
			input.dy = -1 // Move para cima
		case 'a':
			input.player = 0
			input.dx = -1 // Move para a esquerda
		case 's':
			input.player = 0
			input.dy = 1 // Move para baixo
		case 'd':
			input.player = 0
			input.dx = 1 // Move para a direita
		case 'i':
			input.player = 1
			input.dy = -1 // Move para cima
		case 'j':
			input.player = 1
			input.dx = -1 // Move para a esquerda
		case 'k':
			input.player = 1
			input.dy = 1 // Move para baixo
		case 'l':
			input.player = 1
			input.dx = 1
		}

		if input.player == 0 {
			player1Input <- input
		} else {
			player2Input <- input
		}
	}
	return true // Continua o jogo
}