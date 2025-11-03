// main.go - Loop principal do jogo

package main

import (
	"fmt"
	"net/rpc"
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

// Structs para RPC com o servidor
type Posicoes struct {
	Pos1X, Pos1Y, Pos2X, Pos2Y int
	ClientID                   int
	SequenceNumber             int
}

type MoverElementoTypeRPC struct {
	Player         int
	X, Y           int
	ClientID       int
	SequenceNumber int
}

// Função para chamadas RPC com retry automático
func callRPCWithRetry(client *rpc.Client, method string, args interface{}, reply interface{}, maxRetries int) error {
	var err error
	for tentativa := 0; tentativa < maxRetries; tentativa++ {
		err = client.Call(method, args, reply)
		if err == nil {
			return nil
		}
		fmt.Printf("Erro na chamada RPC '%s' (tentativa %d/%d): %v\n", method, tentativa+1, maxRetries, err)
		time.Sleep(time.Duration(tentativa+1) * 100 * time.Millisecond) // Backoff exponencial
	}
	return fmt.Errorf("falha após %d tentativas: %w", maxRetries, err)
}

func main() {

	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"

	// conecta no servidor
	if len(os.Args) != 2 {
		fmt.Print("É necessário informar um ipv4")
		return
	}
	porta := 8973
	addr := os.Args[1]
	fmt.Printf("Conectando ao servidor em %s na porta %d\n", addr, porta)
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", addr, porta))
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	defer client.Close()
	// solicita o número do jogador
	var numeroJogador int
	err = callRPCWithRetry(client, "DadosJogo.ConectarJogador", true, &numeroJogador, 5)
	if err != nil {
		fmt.Println("Erro ao conectar como jogador:", err)
		return
	}
	if numeroJogador == -1 {
		fmt.Println("Servidor cheio. Não foi possível conectar como jogador.")
		return
	}
	fmt.Printf("Conectado como Jogador %d\n", numeroJogador)

	// Inicializa o jogo
	jogo := jogoNovo(numeroJogador, client)
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
				apagarFogo(&jogo)
			}
			// Verifica colisão inimigo de fogo com personagem de água
			if jogo.IniFogoPosX == jogo.Pos2X && jogo.IniFogoPosY == jogo.Pos2Y {
				// Volta personagem de água para posição inicial
				evaporarAgua(&jogo)
			}
			interfaceDesenharJogo(&jogo)
			time.Sleep(16 * time.Millisecond)
		}
	}()

	go func() {
		for {
			// Solicita posições atualizadas dos jogadores ao servidor
			var posicoes Posicoes
			err := callRPCWithRetry(client, "DadosJogo.GetPosicoes", numeroJogador, &posicoes, 3)
			if err != nil {
				fmt.Println("Erro ao obter posições dos jogadores:", err)
				continue
			}
			if numeroJogador == 1 {
				jogo.Pos2X = posicoes.Pos2X
				jogo.Pos2Y = posicoes.Pos2Y
				if jogo.Mapa[jogo.Pos2Y][jogo.Pos2X].simbolo == BandeiraAgua.simbolo {
					// Se eu sou o jogador 1, o jogador 2 venceu
					if jogo.JogadorAtual == 1 {
						jogo.StatusMsg = "JOGADOR 2 (ÁGUA) VENCEU!"
						player2Vence <- true
					}
				}
			} else if numeroJogador == 2 {
				jogo.Pos1X = posicoes.Pos1X
				jogo.Pos1Y = posicoes.Pos1Y
				if jogo.Mapa[jogo.Pos1Y][jogo.Pos1X].simbolo == BandeiraFogo.simbolo {
					// Se eu sou o jogador 2, o jogador 1 venceu
					if jogo.JogadorAtual == 2 {
						jogo.StatusMsg = "JOGADOR 1 (FOGO) VENCEU!"
						player1Vence <- true
					}
				}
			}

			time.Sleep(16 * time.Millisecond)
		}
	}()

	go recebeInput(0, &jogo)
	go recebeInput(1, &jogo)
	go inimigoRecebeInput(0, &jogo)
	go inimigoRecebeInput(1, &jogo)
	go inimigoPatrulha(0, &jogo)
	go inimigoPatrulha(1, &jogo)
	go ativarBotoes(&jogo)
	go jogoMoverElemento()
	go vencerJogo(&jogo)

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
