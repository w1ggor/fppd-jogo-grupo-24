// personagem.go - Funções para movimentação e ações do personagem
package main

import (
	"time"
)

// Atualiza a posição do personagem com base na tecla pressionada (WASD)
func personagemMover(input InputData, jogo *Jogo, player int) {

	dx, dy := input.dx, input.dy

	if player == 0 {
		nx, ny := jogo.Pos1X+dx, jogo.Pos1Y+dy
		// Verifica se o movimento é permitido e realiza a movimentação
		if jogoPodeMoverPara(jogo, nx, ny, player) {
			var moveInput = MoverElementoType{player: 0, jogo: jogo, x: jogo.Pos1X, y: jogo.Pos1Y, dx: dx, dy: dy}
			moveElemento <- moveInput
			jogo.Pos1X, jogo.Pos1Y = nx, ny

			var sucesso bool
			arg := MoverElementoTypeRPC{Player: 1, X: jogo.Pos1X, Y: jogo.Pos1Y}
			jogo.cliente.Call("DadosJogo.MoverJogador", arg, &sucesso)
		}
	} else {
		nx, ny := jogo.Pos2X+dx, jogo.Pos2Y+dy
		// Verifica se o movimento é permitido e realiza a movimentação
		if jogoPodeMoverPara(jogo, nx, ny, player) {
			var moveInput = MoverElementoType{player: 1, jogo: jogo, x: jogo.Pos2X, y: jogo.Pos2Y, dx: dx, dy: dy}
			moveElemento <- moveInput
			jogo.Pos2X, jogo.Pos2Y = nx, ny

			var sucesso bool
			arg := MoverElementoTypeRPC{Player: 2, X: jogo.Pos2X, Y: jogo.Pos2Y}
			jogo.cliente.Call("DadosJogo.MoverJogador", arg, &sucesso)
		}
	}

}

// Define o que ocorre quando o jogador pressiona a tecla de interação
// Neste exemplo, apenas exibe uma mensagem de status
// Você pode expandir essa função para incluir lógica de interação com objetos

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
	case "mover":
		// Move o personagem com base na tecla
		switch ev.Tecla {
		case 'w':
			input.dy = -1 // Move para cima
		case 'a':

			input.dx = -1 // Move para a esquerda
		case 's':

			input.dy = 1 // Move para baixo
		case 'd':

			input.dx = 1 // Move para a direita
		}
		input.player = jogo.JogadorAtual
		if input.player == 1 {
			player1Input <- input
		} else {
			player2Input <- input
		}
	}
	return true // Continua o jogo
}

var player1Vence = make(chan bool, 1)
var player2Vence = make(chan bool, 1)

func vencerJogo(jogo *Jogo) {

	jogador1chegou := false
	jogador2chegou := false
	jogo.StatusMsg = "Voces devem chegar nas bandeiras juntos "
	for !jogador1chegou || !jogador2chegou {

		select {
		case <-player1Vence:
			jogador1chegou = true
			jogo.StatusMsg = "Jogador 1 chegou!"
		case <-player2Vence:
			jogador2chegou = true
			jogo.StatusMsg = "Jogador 2 chegou!"
		}
	}

	jogo.StatusMsg = "Voces Ganharam!!!!"
	time.Sleep(time.Second * 2)
	resetPersonagens(jogo)
	vencerJogo(jogo)
}
func resetPersonagens(jogo *Jogo) {
	jogo.Pos1X, jogo.Pos1Y = jogo.PosCo1X, jogo.PosCo1Y
	jogo.Pos2X, jogo.Pos2Y = jogo.PosCo2X, jogo.PosCo2Y
	jogo.UltimoVisitado1 = Vazio
	jogo.UltimoVisitado2 = Vazio
}

func apagarFogo(jogo *Jogo) {
	// Salva o elemento atual para restaurar depois
	elementoAtual := jogo.Mapa[jogo.Pos1Y][jogo.Pos1X]

	// Move o personagem para a posição inicial
	jogo.Mapa[jogo.Pos1Y][jogo.Pos1X] = jogo.UltimoVisitado1
	jogo.Pos1X, jogo.Pos1Y = jogo.PosCo1X, jogo.PosCo1Y
	jogo.UltimoVisitado1 = jogo.Mapa[jogo.Pos1Y][jogo.Pos1X]
	jogo.Mapa[jogo.Pos1Y][jogo.Pos1X] = elementoAtual

	jogo.StatusMsg = "Fogo apagou!"
}

func evaporarAgua(jogo *Jogo) {
	// Salva o elemento atual para restaurar depois
	elementoAtual := jogo.Mapa[jogo.Pos2Y][jogo.Pos2X]

	// Move o personagem para a posição inicial
	jogo.Mapa[jogo.Pos2Y][jogo.Pos2X] = jogo.UltimoVisitado2
	jogo.Pos2X, jogo.Pos2Y = jogo.PosCo2X, jogo.PosCo2Y
	jogo.UltimoVisitado2 = jogo.Mapa[jogo.Pos2Y][jogo.Pos2X]
	jogo.Mapa[jogo.Pos2Y][jogo.Pos2X] = elementoAtual

	jogo.StatusMsg = "Agua evaporou!"
}
