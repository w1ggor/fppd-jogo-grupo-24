package main

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

type DadosJogo struct {
	player1, player2 bool
	posicaoJogadores Posicoes
	mu               sync.Mutex
}
type Posicoes struct {
	Pos1X, Pos1Y, Pos2X, Pos2Y int
}
type Inicializar struct {
	Pos1X, Pos1Y, Pos2X, Pos2Y int
}

func (s *DadosJogo) Inicializar(dados Inicializar, sucesso *bool) error {
	s.posicaoJogadores.Pos1X = dados.Pos1X
	s.posicaoJogadores.Pos1Y = dados.Pos1Y
	s.posicaoJogadores.Pos2X = dados.Pos2X
	s.posicaoJogadores.Pos2Y = dados.Pos2Y
	*sucesso = true
	return nil
}

func (s *DadosJogo) ConectarJogador(_ bool, numeroJogador *int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.player1 {
		s.player1 = true
		*numeroJogador = 1
		fmt.Println("Jogador 1 conectado")
	} else if !s.player2 {
		s.player2 = true
		*numeroJogador = 2
		fmt.Println("Jogador 2 conectado")
	} else {
		*numeroJogador = -1
		fmt.Println("Tentativa de conexão de um terceiro jogador")
	}
	return nil
}

func (s *DadosJogo) GetPosicoes(player int, resposta *Posicoes) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	*resposta = s.posicaoJogadores
	fmt.Println("Jogador", player, "solicitou posições dos jogadores")
	return nil
}

type MoverElementoType struct {
	Player int
	X, Y   int
}

func (s *DadosJogo) MoverElemento(dados MoverElementoType, sucesso *bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if dados.Player == 1 {
		s.posicaoJogadores.Pos1X = dados.X
		s.posicaoJogadores.Pos1Y = dados.Y
		fmt.Println("Movimentação do Jogador 1 para:", dados.X, dados.Y)
	} else if dados.Player == 2 {
		s.posicaoJogadores.Pos2X = dados.X
		s.posicaoJogadores.Pos2Y = dados.Y
		fmt.Println("Movimentação do Jogador 2 para:", dados.X, dados.Y)
	}
	*sucesso = true
	return nil
}

func main() {
	porta := 8973
	servidor := new(DadosJogo)
	rpc.Register(servidor)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", porta))
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}

	for {
		fmt.Println("Servidor aguardando conexões na porta", porta)
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
