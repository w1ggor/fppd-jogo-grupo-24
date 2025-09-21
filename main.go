// main.go - Loop principal do jogo

package main

import (
	"os"
	"time"
)

// Função auxiliar para valor absoluto
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type InputData struct {
	player int
	input  EventoTeclado
	dx, dy int
}

func main() {

	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// Atualiza a tela periodicamente para mostrar movimentação dos inimigos
	go func() {
		for {
			// Verifica colisão inimigo de água com personagem de fogo
			if jogo.IniAguaPosX == jogo.Pos1X && jogo.IniAguaPosY == jogo.Pos1Y {
				// Volta personagem de fogo para posição inicial
				// Procura posição inicial no mapa
				personagemMover(InputData{player: 0, dx: jogo.PosCo1X - jogo.Pos1X, dy: jogo.PosCo1Y - jogo.Pos1Y}, &jogo, 0) // Atualiza posição na tela

				jogo.StatusMsg = "Fogo apagou!"
			}
			// Verifica colisão inimigo de fogo com personagem de água
			if jogo.IniFogoPosX == jogo.Pos2X && jogo.IniFogoPosY == jogo.Pos2Y {
				// Volta personagem de água para posição inicial
				personagemMover(InputData{player: 1, dx: jogo.PosCo2X - jogo.Pos2X, dy: jogo.PosCo2Y - jogo.Pos2Y}, &jogo, 1) // Atualiza posição na tela

				jogo.StatusMsg = "Agua ferveu!"
			}
			interfaceDesenharJogo(&jogo)
			time.Sleep(50 * time.Millisecond)
		}
	}()

	go recebeInput(0, &jogo)
	go recebeInput(1, &jogo)
	go inimigoRecebeInput(0, &jogo)
	go inimigoRecebeInput(1, &jogo)
	go inimigoPatrulha(0, &jogo)
	go inimigoPatrulha(1, &jogo)
	go jogoMoverElemento()

	// Goroutine para monitorar proximidade e alertar inimigos
	go func() {
		for {
			// Inimigo de água acelera se player de fogo está perto
			distAgua := abs(jogo.IniAguaPosX-jogo.Pos1X) + abs(jogo.IniAguaPosY-jogo.Pos1Y)
			if distAgua <= 15 {
				IniAguaAlerta <- true
			} else {
				IniAguaAlerta <- false
			}
			// Inimigo de fogo acelera se player de água está perto
			distFogo := abs(jogo.IniFogoPosX-jogo.Pos2X) + abs(jogo.IniFogoPosY-jogo.Pos2Y)
			if distFogo <= 15 {
				IniFogoAlerta <- true
			} else {
				IniFogoAlerta <- false
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
			break
		}
		interfaceDesenharJogo(&jogo)
	}
}
