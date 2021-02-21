package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter your username: ")
	var username string
	scanner.Scan()
	username = scanner.Text()

	fmt.Println("Enter your port to serve on: ")
	var myPort string
	scanner.Scan()
	myPort = scanner.Text()

	fmt.Println("Enter peer port to connect to: ")
	var peerPort string
	scanner.Scan()
	peerPort = scanner.Text()

	peer := NewPeer(username, myPort, peerPort)

	go peer.Start()
	fmt.Println("Started")
	for {
		var message string
		fmt.Print("> ")
		scanner.Scan()
		message = scanner.Text()

		if message == "kill" {
			fmt.Println("Kill detected.")
			go peer.Stop()
			return
		}
		peer.SendMessage(message)
	}
}

// func one(wg *sync.WaitGroup, ring *RingBuffer) {
// 	defer wg.Done()

// 	for {
// 		toWrite := make([]byte, 100)
// 		contents := []byte("hello")
// 		copy(toWrite, contents)
// 		ring.Write(toWrite)
// 		time.Sleep(2)
// 	}

// }

// func two(wg *sync.WaitGroup, ring *RingBuffer) {
// 	defer wg.Done()

// 	for {
// 		toRead := make([]byte, 100)
// 		n, err := ring.Read(toRead)
// 		fmt.Println("Read out ", n, "errors", err)
// 		s := string(toRead)
// 		fmt.Println(s)
// 		time.Sleep(1)
// 	}

// }

// func main() {
// 	fmt.Println("Hello World!")
// 	ring := NewRingBuffer(1024)
// 	var wg sync.WaitGroup
// 	go one(&wg, ring)
// 	go two(&wg, ring)
// 	wg.Add(2)
// 	wg.Wait()
// }
