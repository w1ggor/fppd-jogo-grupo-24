// interface.go - Interface gráfica do jogo usando termbox
// O código abaixo implementa a interface gráfica do jogo usando a biblioteca termbox-go.
// A biblioteca termbox-go é uma biblioteca de interface de terminal que permite desenhar
// elementos na tela, capturar eventos do teclado e gerenciar a aparência do terminal.

package main

import (
	"time"

	"github.com/nsf/termbox-go"
)

// Define um tipo Cor para encapsuladar as cores do termbox
type Cor = termbox.Attribute

// Definições de cores utilizadas no jogo
const (
	CorPadrao      Cor = termbox.ColorDefault
	CorCinzaEscuro     = termbox.ColorDarkGray
	CorVermelho        = termbox.ColorRed
	CorAzul            = termbox.ColorBlue
	CorVerde           = termbox.ColorGreen
	CorParede          = termbox.ColorBlack | termbox.AttrBold | termbox.AttrDim
	CorFundoParede     = termbox.ColorDarkGray
	CorTexto           = termbox.ColorDarkGray
)

// EventoTeclado representa uma ação detectada do teclado (como mover, sair ou interagir)
type EventoTeclado struct {
	Tipo  string // "sair", "interagir", "mover"
	Tecla rune   // Tecla pressionada, usada no caso de movimento
}

// Inicializa a interface gráfica usando termbox
func interfaceIniciar() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
}

// Encerra o uso da interface termbox
func interfaceFinalizar() {
	termbox.Close()
}

// Lê um evento do teclado e o traduz para um EventoTeclado
func interfaceLerEventoTeclado() EventoTeclado {
	ev := termbox.PollEvent()
	if ev.Type != termbox.EventKey {
		return EventoTeclado{}
	}
	if ev.Key == termbox.KeyEsc {
		return EventoTeclado{Tipo: "sair"}
	}
	if ev.Ch == 'e' {
		return EventoTeclado{Tipo: "interagir"}
	}
	return EventoTeclado{Tipo: "mover", Tecla: ev.Ch}
}

// Renderiza todo o estado atual do jogo na tela
func interfaceDesenharJogo(jogo *Jogo) {
	interfaceLimparTela()

	// Desenha todos os elementos do mapa
	for y, linha := range jogo.Mapa {
		for x, elem := range linha {
			interfaceDesenharElemento(x, y, elem)
		}
	}

	// Desenha o personagem sobre o mapa
	interfaceDesenharElemento(jogo.Pos1X, jogo.Pos1Y, PersonagemFogo)
	interfaceDesenharElemento(jogo.Pos2X, jogo.Pos2Y, PersonagemAgua)
	// Desenha os inimigos sobre o mapa
	interfaceDesenharElemento(jogo.IniFogoPosX, jogo.IniFogoPosY, InimigoFogo)
	interfaceDesenharElemento(jogo.IniAguaPosX, jogo.IniAguaPosY, InimigoAgua)
	// Desenha a barra de status
	interfaceDesenharBarraDeStatus(jogo)
	// Desenha o portao abrindo
	interfaceDesenharElemento(jogo.PosPortao1XA, jogo.PosPortao1YA, Vazio)
	interfaceDesenharElemento(jogo.PosPortao2XA, jogo.PosPortao2YA, Vazio)
	// Desenha o portao fechando
	interfaceDesenharElemento(jogo.PosPortao1XF, jogo.PosPortao1YF, Portao)
	interfaceDesenharElemento(jogo.PosPortao2XF, jogo.PosPortao2YF, Portao)
	// Força a atualização do terminal
	interfaceAtualizarTela()
	time.Sleep(time.Millisecond * 16)
}

// Limpa a tela do terminal
func interfaceLimparTela() {
	termbox.Clear(CorPadrao, CorPadrao)
}

// Força a atualização da tela do terminal com os dados desenhados
func interfaceAtualizarTela() {
	termbox.Flush()
}

// Desenha um elemento na posição (x, y)
func interfaceDesenharElemento(x, y int, elem Elemento) {
	termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
}

// Exibe uma barra de status com informações úteis ao jogador
func interfaceDesenharBarraDeStatus(jogo *Jogo) {
	// Linha de status dinâmica
	for i, c := range jogo.StatusMsg {
		termbox.SetCell(i, len(jogo.Mapa)+1, c, CorTexto, CorPadrao)
	}

	// Instruções fixas
	msg := "Use WASD para mover o personagem de FOGO"
	for i, c := range msg {
		termbox.SetCell(i, len(jogo.Mapa)+3, c, CorTexto, CorVermelho)
	}

	// Instruções fixas
	msg2 := "Use IJKL para mover o personagem de AGUA."
	for i, c := range msg2 {
		termbox.SetCell(i, len(jogo.Mapa)+4, c, CorTexto, CorAzul)
	}

	// Instruções fixas
	msg3 := "ESC para sair."
	for i, c := range msg3 {
		termbox.SetCell(i, len(jogo.Mapa)+5, c, CorTexto, CorPadrao)
	}
}
