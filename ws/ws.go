package ws

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
)

var ActivePeers = make(map[string]net.Conn)

var InvertedPeers = make(map[net.Conn]string)

var Rooms = make(map[string]Room)

type Room struct {
	Id    int
	Peers []Peer
}

type Peer struct {
	username string
	conn     net.Conn
	online   bool
}

func NewPeer(username string, conn net.Conn) (*Peer, error) {
	usr := strings.TrimSpace(username)
	peer := &Peer{
		username: usr,
		conn:     conn,
		online:   true,
	}
	_, ok := ActivePeers[usr]
	if ok {
		return nil, fmt.Errorf("username already exists")
	}
	ActivePeers[usr] = conn
	_, ok = InvertedPeers[conn]
	if ok {
		return nil, fmt.Errorf("you cant have two usernames")
	}
	InvertedPeers[conn] = usr
	return peer, nil
}

/*
* @param name - represents the name of the room
 */
func (p *Peer) NewRoom(name string) {
	id := int(rand.Intn(100))
	newRoom := &Room{
		Id:    id,
		Peers: []Peer{},
	}
	newRoom.Peers = append(newRoom.Peers, *p)
	Rooms[name] = *newRoom
}

func SendPrivateMessage(sender, receiver string, message []string) error {
	conn, ok := ActivePeers[receiver]
	if !ok {
		return fmt.Errorf("%s is offline", receiver)
	}
	conn.Write([]byte(fmt.Sprintf("(new message) from %s: %s", sender, strings.Join(message, " "))))
	return nil
}
