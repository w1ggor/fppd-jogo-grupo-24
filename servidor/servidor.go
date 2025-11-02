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
	pos1X, pos1Y, pos2X, pos2Y int
}

func (s *DadosJogo) Inicializar(dados Inicializar, sucesso *bool) error {
	s.posicaoJogadores.Pos1X = dados.pos1X
	s.posicaoJogadores.Pos1Y = dados.pos1Y
	s.posicaoJogadores.Pos2X = dados.pos2X
	s.posicaoJogadores.Pos2Y = dados.pos2Y
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

func (s *DadosJogo) GetPosicoes(args struct{}, resposta *Posicoes) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	*resposta = s.posicaoJogadores
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
