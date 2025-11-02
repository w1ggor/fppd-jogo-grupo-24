package main

import (
	"bufio"
	"net/rpc"
	"os"
	"time"
)

type Elemento struct {
	simbolo  rune
	cor      Cor
	corFundo Cor
	tangivel bool
}
type MoverElementoType struct {
	jogo         *Jogo
	player       int
	x, y, dx, dy int
}

type Jogo struct {
	cliente                            *rpc.Client
	JogadorAtual                       int
	sequenceNumber                     int
	Mapa                               [][]Elemento
	PosCo1X, PosCo1Y, PosCo2X, PosCo2Y int
	Pos1X, Pos1Y, Pos2X, Pos2Y         int
	IniFogoPosX, IniFogoPosY           int
	IniAguaPosX, IniAguaPosY           int
	UltimoVisitado1                    Elemento
	UltimoVisitado2                    Elemento
	PosPortao1XF, PosPortao1YF         int
	PosPortao2XF, PosPortao2YF         int
	PosPortao1XA, PosPortao1YA         int
	PosPortao2XA, PosPortao2YA         int
	StatusMsg                          string
	LogMsg                             string
}

func (jogo *Jogo) proximoSequenceNumber() int {
	jogo.sequenceNumber++
	return jogo.sequenceNumber
}

var (
	PersonagemFogo = Elemento{'○', CorVermelho, CorPadrao, true}
	PersonagemAgua = Elemento{'●', CorAzul, CorPadrao, true}
	Inimigo        = Elemento{'☠', CorVermelho, CorPadrao, true}
	Personagem     = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	InimigoFogo    = Elemento{'◇', CorVermelho, CorPadrao, true}
	InimigoAgua    = Elemento{'◆', CorAzul, CorPadrao, true}
	Parede         = Elemento{'▤', CorParede, CorFundoParede, true}
	Portao         = Elemento{'▒', CorPadrao, CorPadrao, true}
	Botao          = Elemento{'◙', CorPadrao, CorPadrao, false}
	Vegetacao      = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio          = Elemento{' ', CorPadrao, CorPadrao, false}
	Fogo           = Elemento{'^', CorVermelho, CorPadrao, false}
	Agua           = Elemento{'~', CorAzul, CorPadrao, false}
	BandeiraFogo   = Elemento{'⚐', CorVermelho, CorPadrao, false}
	BandeiraAgua   = Elemento{'⚑', CorAzul, CorPadrao, false}
)

func jogoNovo(numeroJogador int, cliente *rpc.Client) Jogo {
	return Jogo{
		JogadorAtual:    numeroJogador,
		cliente:         cliente,
		sequenceNumber:  0,
		UltimoVisitado1: Vazio,
		UltimoVisitado2: Vazio,
	}
}

func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case Parede.simbolo:
				e = Parede
			case InimigoFogo.simbolo:
				jogo.IniFogoPosX, jogo.IniFogoPosY = x, y
				e = Vazio
			case InimigoAgua.simbolo:
				jogo.IniAguaPosX, jogo.IniAguaPosY = x, y
				e = Vazio
			case Portao.simbolo:
				e = Portao
			case Botao.simbolo:
				e = Botao
			case Vegetacao.simbolo:
				e = Vegetacao
			case PersonagemFogo.simbolo:
				jogo.PosCo1X, jogo.PosCo1Y = x, y
				jogo.Pos1X, jogo.Pos1Y = x, y
			case PersonagemAgua.simbolo:
				jogo.PosCo2X, jogo.PosCo2Y = x, y
				jogo.Pos2X, jogo.Pos2Y = x, y
			case Fogo.simbolo:
				e = Fogo
			case Agua.simbolo:
				e = Agua
			case BandeiraFogo.simbolo:
				e = BandeiraFogo
			case BandeiraAgua.simbolo:
				e = BandeiraAgua

			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	var sucesso bool
	var pos = Posicoes{
		Pos1X:          jogo.PosCo1X,
		Pos1Y:          jogo.PosCo1Y,
		Pos2X:          jogo.PosCo2X,
		Pos2Y:          jogo.PosCo2Y,
		ClientID:       jogo.JogadorAtual,
		SequenceNumber: jogo.proximoSequenceNumber(),
	}

	erro := ChamarRPC(jogo.cliente, "DadosJogo.Inicializar", pos, &sucesso, 5)
	if erro != nil {
		return erro
	}
	return nil
}

func jogoPodeMoverPara(jogo *Jogo, x, y int, player ...int) bool {
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	if jogo.Mapa[y][x].tangivel {
		return false
	}

	if jogo.Mapa[y][x].simbolo == Agua.simbolo && player != nil && player[0] == 0 {
		apagarFogo(jogo)
		return false
	}
	if jogo.Mapa[y][x].simbolo == Fogo.simbolo && player != nil && player[0] == 1 {
		evaporarAgua(jogo)
		return false
	}
	if jogo.Mapa[y][x].simbolo == BandeiraFogo.simbolo && player != nil && player[0] == 0 {
		jogo.StatusMsg = "O FOGO CHEGOU !"
		player1Vence <- true

		return true
	}
	if jogo.Mapa[y][x].simbolo == BandeiraAgua.simbolo && player != nil && player[0] == 1 {
		jogo.StatusMsg = "A ÁGUA CHEGOU !"
		player2Vence <- true
		return true
	}
	return true
}

var moveElemento = make(chan MoverElementoType, 1)

func jogoMoverElemento() {
	for {
		var moveInput = <-moveElemento
		var jogo = moveInput.jogo
		var player, x, y, dx, dy = moveInput.player, moveInput.x, moveInput.y, moveInput.dx, moveInput.dy
		nx, ny := x+dx, y+dy

		if ny < 0 || ny >= len(jogo.Mapa) || nx < 0 || nx >= len(jogo.Mapa[ny]) {
			continue
		}

		if jogo.Mapa[ny][nx].simbolo == Agua.simbolo || jogo.Mapa[ny][nx].simbolo == Fogo.simbolo {
			continue
		}

		if jogo.Mapa[y][x].simbolo == Agua.simbolo || jogo.Mapa[y][x].simbolo == Fogo.simbolo {
			continue
		}
		if jogo.Mapa[y][x].simbolo == Botao.simbolo || jogo.Mapa[ny][nx].simbolo == Botao.simbolo {
			continue
		}
		if jogo.Mapa[y][x].simbolo == Portao.simbolo || jogo.Mapa[ny][nx].simbolo == Portao.simbolo {
			continue
		}
		if jogo.Mapa[y][x].simbolo == BandeiraAgua.simbolo || jogo.Mapa[ny][nx].simbolo == BandeiraAgua.simbolo {
			continue
		}
		if jogo.Mapa[y][x].simbolo == BandeiraFogo.simbolo || jogo.Mapa[ny][nx].simbolo == BandeiraFogo.simbolo {
			continue
		}

		elemento := jogo.Mapa[y][x]

		if player == 0 {
			jogo.Mapa[y][x] = jogo.UltimoVisitado1
			jogo.UltimoVisitado1 = jogo.Mapa[ny][nx]
			jogo.Mapa[ny][nx] = elemento
		} else if player == 1 {
			jogo.Mapa[y][x] = jogo.UltimoVisitado2
			jogo.UltimoVisitado2 = jogo.Mapa[ny][nx]
			jogo.Mapa[ny][nx] = elemento
		} else {
			jogo.Mapa[y][x] = Vazio
			jogo.Mapa[ny][nx] = elemento
		}
	}

}

func ativarBotoes(jogo *Jogo) {
	var canalP1 = make(chan int)
	var canalP2 = make(chan int)
	var interromperP1 = make(chan int)
	var interromperP2 = make(chan int)

	go ativarB1(jogo, canalP2, interromperP2)
	go ativarB2(jogo, canalP1, interromperP1)
}

func ativarB1(jogo *Jogo, canalP2 chan int, interromperP2 chan int) {
	for {
		if jogo.Pos2X == 13 && jogo.Pos2Y == 12 {
			go abrirP2(jogo, canalP2)
			<-canalP2
			go fecharP2(jogo, interromperP2)
			<-interromperP2
		}
	}
}

func ativarB2(jogo *Jogo, canalP1 chan int, interromperP1 chan int) {
	for {
		if jogo.Pos1X == 66 && jogo.Pos1Y == 24 {
			go abrirP1(jogo, canalP1)
			<-canalP1
			go fecharP1(jogo, interromperP1)
			<-interromperP1
		}
	}
}

func abrirP1(jogo *Jogo, canalP1 chan int) {
	px1 := 25
	py1 := 17

	for {
		if px1 == 0 {
			canalP1 <- 1
			return
		}
		jogo.Mapa[py1][px1] = Vazio
		jogo.PosPortao1XA, jogo.PosPortao1YA = px1, py1
		time.Sleep(time.Millisecond * 100)
		px1--
	}
}

func abrirP2(jogo *Jogo, canalP2 chan int) {
	px2 := 78
	py2 := 17

	for {
		if px2 == 53 {
			canalP2 <- 53
			return
		}
		jogo.Mapa[py2][px2] = Vazio
		jogo.PosPortao2XA, jogo.PosPortao2YA = px2, py2
		time.Sleep(time.Millisecond * 100)
		px2--
	}
}

func fecharP1(jogo *Jogo, interromperP1 chan int) {
	px1 := 1
	py1 := 17

	for {
		if jogo.Pos1X == 66 && jogo.Pos1Y == 24 {
			continue
		}
		for px1 < 26 {
			jogo.Mapa[py1][px1] = Portao
			jogo.PosPortao1XF, jogo.PosPortao1YF = px1, py1
			time.Sleep(time.Millisecond * 100)
			px1++
		}
		interromperP1 <- 25
		return
	}

}

func fecharP2(jogo *Jogo, interromperP2 chan int) {
	px2 := 54
	py2 := 17

	for {
		if jogo.Pos2X == 13 && jogo.Pos2Y == 12 {
			continue
		}
		for px2 < 79 {
			jogo.Mapa[py2][px2] = Portao
			jogo.PosPortao2XF, jogo.PosPortao2YF = px2, py2
			time.Sleep(time.Millisecond * 100)
			px2++
		}
		interromperP2 <- 78
		return
	}
}
