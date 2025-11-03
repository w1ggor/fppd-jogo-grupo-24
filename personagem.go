package main

import (
	"time"
)

func personagemMover(input InputData, jogo *Jogo, player int) {

	dx, dy := input.dx, input.dy

	if player == 0 {
		nx, ny := jogo.Pos1X+dx, jogo.Pos1Y+dy
		if jogoPodeMoverPara(jogo, nx, ny, player) {
			var moveInput = MoverElementoType{player: 0, jogo: jogo, x: jogo.Pos1X, y: jogo.Pos1Y, dx: dx, dy: dy}
			moveElemento <- moveInput
			jogo.Pos1X, jogo.Pos1Y = nx, ny

			var sucesso bool
			arg := MoverElementoTypeRPC{
				Player:         1,
				X:              jogo.Pos1X,
				Y:              jogo.Pos1Y,
				ClientID:       jogo.JogadorAtual,
				SequenceNumber: jogo.proximoSequenceNumber(),
			}
			err := ChamarRPC(jogo.cliente, "DadosJogo.MoverJogador", arg, &sucesso, 3)
			if err != nil {
			}
		}
	} else {
		nx, ny := jogo.Pos2X+dx, jogo.Pos2Y+dy
		if jogoPodeMoverPara(jogo, nx, ny, player) {
			var moveInput = MoverElementoType{player: 1, jogo: jogo, x: jogo.Pos2X, y: jogo.Pos2Y, dx: dx, dy: dy}
			moveElemento <- moveInput
			jogo.Pos2X, jogo.Pos2Y = nx, ny

			var sucesso bool
			arg := MoverElementoTypeRPC{
				Player:         2,
				X:              jogo.Pos2X,
				Y:              jogo.Pos2Y,
				ClientID:       jogo.JogadorAtual,
				SequenceNumber: jogo.proximoSequenceNumber(),
			}
			err := ChamarRPC(jogo.cliente, "DadosJogo.MoverJogador", arg, &sucesso, 3)
			if err != nil {
			}
		}
	}

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

func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo) bool {
	var input = InputData{player: 0, input: ev, dx: 0, dy: 0}
	switch ev.Tipo {
	case "sair":
		return false
	case "mover":
		switch ev.Tecla {
		case 'w':
			input.dy = -1
		case 'a':

			input.dx = -1
		case 's':

			input.dy = 1
		case 'd':

			input.dx = 1
		}
		input.player = jogo.JogadorAtual
		if input.player == 1 {
			player1Input <- input
		} else {
			player2Input <- input
		}
	}
	return true
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
	elementoAtual := jogo.Mapa[jogo.Pos1Y][jogo.Pos1X]

	jogo.Mapa[jogo.Pos1Y][jogo.Pos1X] = jogo.UltimoVisitado1
	jogo.Pos1X, jogo.Pos1Y = jogo.PosCo1X, jogo.PosCo1Y
	jogo.UltimoVisitado1 = jogo.Mapa[jogo.Pos1Y][jogo.Pos1X]
	jogo.Mapa[jogo.Pos1Y][jogo.Pos1X] = elementoAtual

	jogo.StatusMsg = "Fogo apagou!"
}

func evaporarAgua(jogo *Jogo) {
	elementoAtual := jogo.Mapa[jogo.Pos2Y][jogo.Pos2X]

	jogo.Mapa[jogo.Pos2Y][jogo.Pos2X] = jogo.UltimoVisitado2
	jogo.Pos2X, jogo.Pos2Y = jogo.PosCo2X, jogo.PosCo2Y
	jogo.UltimoVisitado2 = jogo.Mapa[jogo.Pos2Y][jogo.Pos2X]
	jogo.Mapa[jogo.Pos2Y][jogo.Pos2X] = elementoAtual

	jogo.StatusMsg = "Agua evaporou!"
}
