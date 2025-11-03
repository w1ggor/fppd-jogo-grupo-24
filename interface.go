package main

import (
	"time"

	"github.com/nsf/termbox-go"
)

type Cor = termbox.Attribute

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

type EventoTeclado struct {
	Tipo  string
	Tecla rune
}

func interfaceIniciar() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
}

func interfaceFinalizar() {
	termbox.Close()
}

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

func interfaceDesenharJogo(jogo *Jogo) {
	interfaceLimparTela()

	for y, linha := range jogo.Mapa {
		for x, elem := range linha {
			interfaceDesenharElemento(x, y, elem)
		}
	}

	interfaceDesenharElemento(jogo.Pos1X, jogo.Pos1Y, PersonagemFogo)
	interfaceDesenharElemento(jogo.Pos2X, jogo.Pos2Y, PersonagemAgua)
	interfaceDesenharElemento(jogo.IniFogoPosX, jogo.IniFogoPosY, InimigoFogo)
	interfaceDesenharElemento(jogo.IniAguaPosX, jogo.IniAguaPosY, InimigoAgua)
	interfaceDesenharBarraDeStatus(jogo)
	interfaceDesenharElemento(jogo.PosPortao1XA, jogo.PosPortao1YA, Vazio)
	interfaceDesenharElemento(jogo.PosPortao2XA, jogo.PosPortao2YA, Vazio)
	interfaceDesenharElemento(jogo.PosPortao1XF, jogo.PosPortao1YF, Portao)
	interfaceDesenharElemento(jogo.PosPortao2XF, jogo.PosPortao2YF, Portao)
	interfaceAtualizarTela()
	time.Sleep(time.Millisecond * 16)
}

func interfaceLimparTela() {
	termbox.Clear(CorPadrao, CorPadrao)
}

func interfaceAtualizarTela() {
	termbox.Flush()
}

func interfaceDesenharElemento(x, y int, elem Elemento) {
	termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
}

func interfaceDesenharBarraDeStatus(jogo *Jogo) {
	for i, c := range jogo.StatusMsg {
		termbox.SetCell(i, len(jogo.Mapa)+1, c, CorTexto, CorPadrao)
	}

	msg := "Use WASD para mover o personagem de FOGO"
	for i, c := range msg {
		termbox.SetCell(i, len(jogo.Mapa)+3, c, CorTexto, CorVermelho)
	}

	msg2 := "Use IJKL para mover o personagem de AGUA."
	for i, c := range msg2 {
		termbox.SetCell(i, len(jogo.Mapa)+4, c, CorTexto, CorAzul)
	}

	msg3 := "ESC para sair."
	for i, c := range msg3 {
		termbox.SetCell(i, len(jogo.Mapa)+5, c, CorTexto, CorPadrao)
	}
}
