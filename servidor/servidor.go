package main

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

// Estrutura para cache de comandos já processados (exactly-once)
type ComandoProcessado struct {
	SequenceNumber int
	Resultado      bool
}

type DadosJogo struct {
	player1, player2 bool
	posicaoJogadores Posicoes
	mu               sync.Mutex
	// Cache de comandos processados por cliente (ClientID -> último SequenceNumber)
	comandosProcessados map[int]map[int]*ComandoProcessado
}
type Posicoes struct {
	Pos1X, Pos1Y, Pos2X, Pos2Y int
}

type Inicializar struct {
	Pos1X, Pos1Y, Pos2X, Pos2Y int
	ClientID                   int
	SequenceNumber             int
}

func (s *DadosJogo) Inicializar(dados Inicializar, sucesso *bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Verifica se comando já foi processado (exactly-once)
	if s.comandosProcessados[dados.ClientID] != nil {
		if cmd, existe := s.comandosProcessados[dados.ClientID][dados.SequenceNumber]; existe {
			*sucesso = cmd.Resultado
			fmt.Printf("Comando duplicado detectado - ClientID: %d, SeqNum: %d (ignorado)\n", dados.ClientID, dados.SequenceNumber)
			return nil
		}
	} else {
		s.comandosProcessados[dados.ClientID] = make(map[int]*ComandoProcessado)
	}

	// Processa o comando
	s.posicaoJogadores.Pos1X = dados.Pos1X
	s.posicaoJogadores.Pos1Y = dados.Pos1Y
	s.posicaoJogadores.Pos2X = dados.Pos2X
	s.posicaoJogadores.Pos2Y = dados.Pos2Y
	*sucesso = true

	// Registra comando como processado
	s.comandosProcessados[dados.ClientID][dados.SequenceNumber] = &ComandoProcessado{
		SequenceNumber: dados.SequenceNumber,
		Resultado:      true,
	}

	fmt.Printf("Inicialização processada - ClientID: %d, SeqNum: %d\n", dados.ClientID, dados.SequenceNumber)
	return nil
}

func (s DadosJogo) Disconnect(clientId int, sucesso *bool) error {
	if clientId == 1 {
		s.player1 = false
		fmt.Println("Jogador 1 desconectado")
	} else if clientId == 2 {
		s.player2 = false
		fmt.Println("Jogador 2 desconectado")
	}
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
	Player         int
	X, Y           int
	ClientID       int // ID único do cliente (1 ou 2)
	SequenceNumber int // Número de sequência do comando
}

func (s *DadosJogo) MoverJogador(dados MoverElementoType, sucesso *bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Verifica se comando já foi processado (exactly-once)
	if s.comandosProcessados[dados.ClientID] != nil {
		if cmd, existe := s.comandosProcessados[dados.ClientID][dados.SequenceNumber]; existe {
			*sucesso = cmd.Resultado
			fmt.Printf("Movimento duplicado detectado - ClientID: %d, SeqNum: %d (ignorado)\n", dados.ClientID, dados.SequenceNumber)
			return nil
		}
	} else {
		s.comandosProcessados[dados.ClientID] = make(map[int]*ComandoProcessado)
	}

	// Processa o comando
	if dados.Player == 1 {
		s.posicaoJogadores.Pos1X = dados.X
		s.posicaoJogadores.Pos1Y = dados.Y
		fmt.Printf("Movimentação do Jogador 1 para: (%d, %d) - SeqNum: %d\n", dados.X, dados.Y, dados.SequenceNumber)
	} else if dados.Player == 2 {
		s.posicaoJogadores.Pos2X = dados.X
		s.posicaoJogadores.Pos2Y = dados.Y
		fmt.Printf("Movimentação do Jogador 2 para: (%d, %d) - SeqNum: %d\n", dados.X, dados.Y, dados.SequenceNumber)
	}
	*sucesso = true

	// Registra comando como processado
	s.comandosProcessados[dados.ClientID][dados.SequenceNumber] = &ComandoProcessado{
		SequenceNumber: dados.SequenceNumber,
		Resultado:      true,
	}

	return nil
}

func main() {
	porta := 8973
	servidor := new(DadosJogo)
	servidor.comandosProcessados = make(map[int]map[int]*ComandoProcessado)
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
