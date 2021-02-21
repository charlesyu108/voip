package main

import (
	"fmt"
	"net"
)

type Peer struct {
	conn     *net.UDPConn
	peerAddr *net.UDPAddr
	myAddr   *net.UDPAddr
	myPort   string
	peerPort string
	name     string
	outbox   chan string
	kill     chan bool
}

func NewPeer(name string, myPort string, peerPort string) *Peer {
	peer := new(Peer)
	peer.name = name
	peer.myPort = myPort
	peer.peerPort = peerPort
	peer.kill = make(chan bool)
	peer.outbox = make(chan string)
	return peer
}

func (peer *Peer) Start() {
	// Do connection stuff here.
	myAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:"+peer.myPort)
	peerAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:"+peer.peerPort)

	conn, _ := net.ListenUDP("udp", peerAddr)

	peer.conn = conn
	peer.myAddr = myAddr
	peer.peerAddr = peerAddr

	go peer.receiveMessage()

	for {
		select {
		case <-peer.kill:
			return
		case messageToRelay := <-peer.outbox:
			go peer.relayMessage(messageToRelay)
		default:
			// TODO
		}
	}
}

func (peer *Peer) receiveMessage() {
	receiver := make([]byte, 512)
	for {
		peer.conn.ReadFromUDP(receiver)
		fmt.Println(string(receiver))
	}

}

func (peer *Peer) relayMessage(message string) {
	sendBuf := make([]byte, 512)
	copy(sendBuf, []byte(message))
	_, err := peer.conn.WriteToUDP(sendBuf, peer.myAddr)
	if err != nil {
		fmt.Println("error!", err)
	}
}

func (peer *Peer) Stop() {
	peer.kill <- true
}

func (peer *Peer) SendMessage(message string) {
	peer.outbox <- fmt.Sprint(peer.name, " [@", peer.myPort, "]: ", message)
}
