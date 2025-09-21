// inimigo.go - Funções para movimentação e ações do inimigo

package main

import "time"

func inimigoMover(input InputData, jogo *Jogo, inimigo int) {

	if inimigo == 0 {
		fdx, fdy := input.dx, input.dy
		nx, ny := jogo.IniFogoPosX+fdx, jogo.IniFogoPosY+fdy
		// Verifica se o movimento é permitido e realiza a movimentação
		if jogoPodeMoverPara(jogo, nx, ny) {
			var moveInput = MoverElementoType{jogo: jogo, x: jogo.IniFogoPosX, y: jogo.IniFogoPosY, dx: fdx, dy: fdy}
			moveElemento <- moveInput
			jogo.IniFogoPosX, jogo.IniFogoPosY = nx, ny
		}
	} else {
		adx, ady := input.dx, input.dy
		nx, ny := jogo.IniAguaPosX+adx, jogo.IniAguaPosY+ady
		// Verifica se o movimento é permitido e realiza a movimentação
		if jogoPodeMoverPara(jogo, nx, ny) {
			var moveInput = MoverElementoType{jogo: jogo, x: jogo.IniAguaPosX, y: jogo.IniAguaPosY, dx: adx, dy: ady}
			moveElemento <- moveInput
			jogo.IniAguaPosX, jogo.IniAguaPosY = nx, ny
		}
	}

}

// Canal para patrulha automática
var IniFogoPatrulha = make(chan InputData)
var IniAguaPatrulha = make(chan InputData)

// Canal de alerta para aumentar velocidade
var IniFogoAlerta = make(chan bool)
var IniAguaAlerta = make(chan bool)

// Goroutine do inimigo: escuta patrulha automática e comandos externos
func inimigoRecebeInput(player int, jogo *Jogo) {
	var patrulhaChan chan InputData
	if player == 0 {
		patrulhaChan = IniFogoPatrulha
	} else {
		patrulhaChan = IniAguaPatrulha
	}
	for {
		input := <-patrulhaChan
		// Movimento automático de patrulha
		inimigoMover(input, jogo, player)

	}
}

// Alterna entre patrulha e alerta
func inimigoPatrulha(player int, jogo *Jogo) {
	var patrulhaChan chan InputData
	var alertaChan chan bool
	if player == 0 {
		patrulhaChan = IniFogoPatrulha
		alertaChan = IniFogoAlerta
	} else {
		patrulhaChan = IniAguaPatrulha
		alertaChan = IniAguaAlerta
	}
	dx := 1           // cada inimigo tem sua própria direção
	velocidade := 500 // ms
	emAlerta := false
	for {
		select {
		case emAlerta = <-alertaChan:
			if emAlerta {
				velocidade = 35 // mais rápido
			} else {
				velocidade = 500 // normal
			}
		default:
			// segue lógica normal
		}
		var nx, ny int
		if player == 0 {
			nx, ny = jogo.IniFogoPosX+dx, jogo.IniFogoPosY
		} else {
			nx, ny = jogo.IniAguaPosX+dx, jogo.IniAguaPosY
		}
		input := InputData{player: player, dx: dx, dy: 0}
		patrulhaChan <- input
		sleepMs(velocidade)
		if !jogoPodeMoverPara(jogo, nx, ny) {
			dx = -dx
		}
	}
}

// Função utilitária para dormir em milissegundos
func sleepMs(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
