package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
)

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

	peer := NewPeer(listeningPort, peerAddr)
	go peer.Start()
	fmt.Println("Started! Ctrl-C to quit.")

	killChan := make(chan os.Signal)
	signal.Notify(killChan, os.Interrupt, os.Kill)

	for {
		select {
		case <-killChan:
			fmt.Println("Cleaning up...")
			peer.Stop()
			fmt.Println("Done!")
			os.Exit(0)
		}
	}
}
