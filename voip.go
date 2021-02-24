package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"github.com/gordonklaus/portaudio"
)

type VoipPeer struct {
	conn      *net.UDPConn
	peerAddr  *net.UDPAddr
	localAddr *net.UDPAddr

	audioStream *portaudio.Stream

	recvData chan []byte
	sendData chan []byte
	kill     chan bool
}

func NewVoipPeer(listeningPort string, peerAddr string) *VoipPeer {
	var err error
	peer := new(VoipPeer)
	peer.localAddr, err = net.ResolveUDPAddr("udp", ":"+listeningPort)
	chk(err)
	peer.peerAddr, err = net.ResolveUDPAddr("udp", peerAddr)
	chk(err)
	peer.kill = make(chan bool)
	peer.recvData = make(chan []byte)
	peer.sendData = make(chan []byte)
	portaudio.Initialize()
	peer.audioStream, err = portaudio.OpenDefaultStream(1, 1, 44100, 1024, peer.streamCallback)
	chk(err)

	return peer
}

func (peer *VoipPeer) Start() {
	defer peer.cleanup()

	conn, _ := net.ListenUDP("udp", peer.localAddr)
	peer.conn = conn

	go peer.sendPackets()
	go peer.receivePackets()
	peer.audioStream.Start()

	for {
		select {
		case <-peer.kill:
			return
		default:
			// TODO
		}
	}
}

func (peer *VoipPeer) receivePackets() {
	for {
		receiver := make([]byte, 8096)
		peer.conn.ReadFromUDP(receiver)
		peer.recvData <- receiver
	}
}

func (peer *VoipPeer) sendPackets() {
	for {
		sendBuf := make([]byte, 8096)
		message := <-peer.sendData

		copy(sendBuf, message)
		peer.conn.WriteToUDP(sendBuf, peer.peerAddr)
	}
}

func (peer *VoipPeer) streamCallback(
	in []int32,
	out []int32,
) {
	go func() {
		inputBuf := new(bytes.Buffer)
		binary.Write(inputBuf, binary.BigEndian, in)
		peer.sendData <- inputBuf.Bytes()
	}()

	go func() {
		output := <-peer.recvData
		outputBuf := bytes.NewReader(output)
		err := binary.Read(outputBuf, binary.BigEndian, out)
		chk(err)
	}()
}

func (peer *VoipPeer) cleanup() {
	peer.audioStream.Close()
	peer.conn.Close()
}

func (peer *VoipPeer) Stop() {
	peer.kill <- true
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter listening port: ")
	var listeningPort string
	scanner.Scan()
	listeningPort = scanner.Text()

	fmt.Println("Enter peer address to connect to: ")
	var peerAddr string
	scanner.Scan()
	peerAddr = scanner.Text()

	peer := NewVoipPeer(listeningPort, peerAddr)
	peer.Start()
	fmt.Println("Started")
}

func chk(err error) {
	if err != nil {
		fmt.Print(err)
	}
}
