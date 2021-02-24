package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/gordonklaus/portaudio"
)

// Peer defines a struct representing a voice-over-IP
// streaming peer over a local area network.
type Peer struct {
	conn      *net.UDPConn
	peerAddr  *net.UDPAddr
	localAddr *net.UDPAddr

	audioStream *portaudio.Stream

	recvBuffer           *bytes.Buffer
	sendBuffer           *bytes.Buffer
	playbackBufferFrames int

	quit chan bool
}

// NewPeer instantiates a new Peer.
func NewPeer(listeningPort string, peerAddr string) *Peer {
	peer := new(Peer)
	peer.recvBuffer = new(bytes.Buffer)
	peer.sendBuffer = new(bytes.Buffer)
	// 200 (ms) x 44100 (frames/s) x 1/1000 (s/ms) = 8820 (frames)
	peer.playbackBufferFrames = 8820
	peer.quit = make(chan bool)

	var err error
	peer.localAddr, err = net.ResolveUDPAddr("udp", ":"+listeningPort)
	chk(err)
	peer.peerAddr, err = net.ResolveUDPAddr("udp", peerAddr)
	chk(err)

	portaudio.Initialize()
	peer.audioStream, err = portaudio.OpenDefaultStream(1, 1, 44100, 1024, peer.streamCallback)
	chk(err)

	return peer
}

// Start begins audio streaming for the Peer.
func (peer *Peer) Start() {
	defer peer.cleanup()

	conn, _ := net.ListenUDP("udp", peer.localAddr)
	peer.conn = conn

	go peer.sendPackets()
	go peer.receivePackets()
	peer.audioStream.Start()

	for {
		select {
		case <-peer.quit:
			return
		}
	}
}

// Stop streaming.
func (peer *Peer) Stop() {
	peer.quit <- true
}

// Clean up tasks.
func (peer *Peer) cleanup() {
	peer.audioStream.Close()
	peer.conn.Close()
}

// Background routine to listen and receive network
// packets from the connected peer.
func (peer *Peer) receivePackets() {
	for {
		receiver := make([]byte, 1024)
		peer.conn.ReadFromUDP(receiver)
		if len(receiver) > 0 {
			peer.recvBuffer.Write(receiver)
		}
	}
}

// Background routine for sending ready packets
// to the connected peer.
func (peer *Peer) sendPackets() {
	for {
		if peer.sendBuffer.Len() > 0 {
			outgoing := make([]byte, 1024)
			peer.sendBuffer.Read(outgoing)
			peer.conn.WriteToUDP(outgoing, peer.peerAddr)
		}
	}
}

// Callback function for the PortAudio in/out streams.
// Launches two go-routines to handle.
func (peer *Peer) streamCallback(
	in []int32,
	out []int32,
) {
	// Transfers the incoming audio to the outbound sendBuffer
	go func() {
		binary.Write(peer.sendBuffer, binary.BigEndian, in)
	}()

	// Reads any ready data out from the recvBuffer into the outgoing audio stream.
	go func() {
		if peer.recvBuffer.Len() < peer.playbackBufferFrames*4 {
			return
		}
		toWrite := make([]int32, 1024)
		readOut := make([]byte, 4096)
		peer.recvBuffer.Read(readOut)
		outputBuf := bytes.NewReader(readOut)
		chk(binary.Read(outputBuf, binary.BigEndian, toWrite))
		copy(out, toWrite)
	}()
}

func chk(err error) {
	if err != nil {
		fmt.Print(err)
	}
}
