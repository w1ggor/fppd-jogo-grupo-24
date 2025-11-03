package main

import "time"

func inimigoMover(input InputData, jogo *Jogo, inimigo int) {

	if inimigo == 0 {
		fdx, fdy := input.dx, input.dy
		nx, ny := jogo.IniFogoPosX+fdx, jogo.IniFogoPosY+fdy
		if jogoPodeMoverPara(jogo, nx, ny) {
			var moveInput = MoverElementoType{player: 4, jogo: jogo, x: jogo.IniFogoPosX, y: jogo.IniFogoPosY, dx: fdx, dy: fdy}
			moveElemento <- moveInput
			jogo.IniFogoPosX, jogo.IniFogoPosY = nx, ny
		}
	} else {
		adx, ady := input.dx, input.dy
		nx, ny := jogo.IniAguaPosX+adx, jogo.IniAguaPosY+ady
		if jogoPodeMoverPara(jogo, nx, ny) {
			var moveInput = MoverElementoType{player: 4, jogo: jogo, x: jogo.IniAguaPosX, y: jogo.IniAguaPosY, dx: adx, dy: ady}
			moveElemento <- moveInput
			jogo.IniAguaPosX, jogo.IniAguaPosY = nx, ny
		}
	}

}

var IniFogoPatrulha = make(chan InputData)
var IniAguaPatrulha = make(chan InputData)

var IniFogoAlerta = make(chan bool)
var IniAguaAlerta = make(chan bool)

func inimigoRecebeInput(player int, jogo *Jogo) {
	var patrulhaChan chan InputData
	if player == 0 {
		patrulhaChan = IniFogoPatrulha
	} else {
		patrulhaChan = IniAguaPatrulha
	}
	for {
		input := <-patrulhaChan
		inimigoMover(input, jogo, player)

	}
}

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
	dx := 1
	velocidade := 500
	emAlerta := false
	for {
		select {
		case emAlerta = <-alertaChan:
			if emAlerta {
				velocidade = 35
			} else {
				velocidade = 500
			}
		default:
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

func sleepMs(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
