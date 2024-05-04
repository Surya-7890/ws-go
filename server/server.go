package server

import (
	"fmt"
	"io"
	"net"

	"github.com/Surya7890/ws-go/ws"
)

type Message struct {
	Data   []byte
	Sender string
}

type Server struct {
	address  string
	listener net.Listener
	exit     chan struct{}
	Msg      chan *Message
}

func NewServer(address string) *Server {
	return &Server{
		address: address,
		exit:    make(chan struct{}),
		Msg:     make(chan *Message, 10),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	s.listener = listener
	defer s.listener.Close()

	s.acceptLoop()
	<-s.exit
	close(s.Msg)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("error while accepting connection:", err)
			conn.Close()
			continue
		}
		buffer := make([]byte, 512)
		go func() {
			for {
				conn.Write([]byte("\nEnter Username: "))
				length, err := conn.Read(buffer)
				if err != nil {
					conn.Close()
				} else {
					_, err := ws.NewPeer(string(buffer[:length])+"\n", conn)
					if err != nil {
						conn.Write([]byte(err.Error()))
						continue
					} else {
						go s.readLoop(conn)
						break
					}
				}
			}
		}()
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		length, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("error while reading loop:", err)
			if err == io.EOF {
				conn.Close()
			} else {
				return
			}
		}
		s.Msg <- &Message{
			Data:   buffer[:length],
			Sender: ws.InvertedPeers[conn],
		}
	}
}
