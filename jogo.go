// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"os"
	"time"
)

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	simbolo  rune
	cor      Cor
	corFundo Cor
	tangivel bool // Indica se o elemento bloqueia passagem
}
type MoverElementoType struct {
	jogo         *Jogo
	player       int
	x, y, dx, dy int
}

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa                               [][]Elemento // grade 2D representando o mapa
	PosCo1X, PosCo1Y, PosCo2X, PosCo2Y int          // posição do comeco do personagem
	Pos1X, Pos1Y, Pos2X, Pos2Y         int          // posição atual do personagem
	IniFogoPosX, IniFogoPosY           int          // posição atual do inimigo de fogo
	IniAguaPosX, IniAguaPosY           int          // posição atual do inimigo de fogo
	UltimoVisitado1                    Elemento     // elemento que estava na posição do personagem antes de mover
	UltimoVisitado2                    Elemento
	PosPortao1XF, PosPortao1YF         int
	PosPortao2XF, PosPortao2YF         int
	PosPortao1XA, PosPortao1YA         int
	PosPortao2XA, PosPortao2YA         int
	StatusMsg                          string // mensagem para a barra de status
}

// Elementos visuais do jogo
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
)

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	// O ultimo elemento visitado é inicializado como vazio
	// pois o jogo começa com o personagem em uma posição vazia
	return Jogo{UltimoVisitado1: Vazio, UltimoVisitado2: Vazio}
}

// Lê um arquivo texto linha por linha e constrói o mapa do jogo
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
				jogo.IniFogoPosX, jogo.IniFogoPosY = x, y // registra a posição inicial do inimigo de fogo
				e = Vazio                                 // remove o símbolo do inimigo do mapa
			case InimigoAgua.simbolo:
				jogo.IniAguaPosX, jogo.IniAguaPosY = x, y // registra a posição inicial do inimigo de água
				e = Vazio                                 // remove o símbolo do inimigo do mapa
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
				jogo.Pos2X, jogo.Pos2Y = x, y // registra a posição inicial do personagem
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// Verifica se o personagem pode se mover para a posição (x, y)
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	// Verifica se a coordenada Y está dentro dos limites verticais do mapa
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	// Verifica se a coordenada X está dentro dos limites horizontais do mapa
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	// Verifica se o elemento de destino é tangível (bloqueia passagem)
	if jogo.Mapa[y][x].tangivel {
		return false
	}

	// Pode mover para a posição
	return true
}

var moveElemento = make(chan MoverElementoType, 1)

// Move um elemento para a nova posição
func jogoMoverElemento() {
	for {
		var moveInput = <-moveElemento
		var jogo = moveInput.jogo
		var player, x, y, dx, dy = moveInput.player, moveInput.x, moveInput.y, moveInput.dx, moveInput.dy
		nx, ny := x+dx, y+dy

		// Obtem elemento atual na posição
		elemento := jogo.Mapa[y][x] // guarda o conteúdo atual da posição
		switch player {
		case 0:
			jogo.Mapa[y][x] = jogo.UltimoVisitado1   // restaura o conteúdo anterior
			jogo.UltimoVisitado1 = jogo.Mapa[ny][nx] // guarda o conteúdo atual da nova posição
			jogo.Mapa[ny][nx] = elemento
		case 1:
			jogo.Mapa[y][x] = jogo.UltimoVisitado2   // restaura o conteúdo anterior
			jogo.UltimoVisitado2 = jogo.Mapa[ny][nx] // guarda o conteúdo atual da nova posição
			jogo.Mapa[ny][nx] = elemento
		default:
			jogo.Mapa[y][x] = Vazio                  // restaura o conteúdo anterior
			jogo.UltimoVisitado2 = jogo.Mapa[ny][nx] // guarda o conteúdo atual da nova posição
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
		if jogo.Pos1X == 13 && jogo.Pos1Y == 12 {
			go abrirP2(jogo, canalP2)
			<-canalP2
			go fecharP2(jogo, interromperP2)
			<-interromperP2
		}
	}
}

func ativarB2(jogo *Jogo, canalP1 chan int, interromperP1 chan int) {
	for {
		if jogo.Pos2X == 66 && jogo.Pos2Y == 24 {
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
		if jogo.Pos2X == 66 && jogo.Pos2Y == 24 {
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
		if jogo.Pos1X == 13 && jogo.Pos1Y == 12 {
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
