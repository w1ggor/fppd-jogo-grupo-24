// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"os"
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
	Vegetacao      = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio          = Elemento{' ', CorPadrao, CorPadrao, false}
	Fogo           = Elemento{'^', CorVermelho, CorPadrao, false}
	Agua           = Elemento{'~', CorAzul, CorPadrao, false}
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
			case Vegetacao.simbolo:
				e = Vegetacao
			case PersonagemFogo.simbolo:
				jogo.PosCo1X, jogo.PosCo1Y = x, y
				jogo.Pos1X, jogo.Pos1Y = x, y
			case PersonagemAgua.simbolo:
				jogo.PosCo2X, jogo.PosCo2Y = x, y
				jogo.Pos2X, jogo.Pos2Y = x, y // registra a posição inicial do personagem                             // remove o símbolo do inimigo do mapa
			case Fogo.simbolo:
				e = Fogo
			case Agua.simbolo:
				e = Agua
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
func jogoPodeMoverPara(jogo *Jogo, x, y int, player ...int) bool {
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

	// Verifica se o elemento de destino é tangível (bloqueia passagem)
	if jogo.Mapa[y][x].simbolo == Agua.simbolo && player != nil && player[0] == 0 {
		apagarFogo(jogo)
		return false
	}
	if jogo.Mapa[y][x].simbolo == Fogo.simbolo && player != nil && player[0] == 1 {
		evaporarAgua(jogo)
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
		if player == 0 {
			jogo.Mapa[y][x] = jogo.UltimoVisitado1   // restaura o conteúdo anterior
			jogo.UltimoVisitado1 = jogo.Mapa[ny][nx] // guarda o conteúdo atual da nova posição
			jogo.Mapa[ny][nx] = elemento
		} else {
			jogo.Mapa[y][x] = jogo.UltimoVisitado2   // restaura o conteúdo anterior
			jogo.UltimoVisitado2 = jogo.Mapa[ny][nx] // guarda o conteúdo atual da nova posição
			jogo.Mapa[ny][nx] = elemento
		}

	}

}
