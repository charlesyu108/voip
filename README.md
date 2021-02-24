# voip
Simple p2p Voice-Over-IP chat over a local area network.

Theoretically, this should be cross-platform supported but only tested on MacOS machines.

# Usage

Build and run:
```
go build && ./voip
```

Install:
```
go install
```
---

The interface will prompt you to provide a listening port and
a peer address to connect to. 

The listening port should be a free port
on your host machine `(e.g. "8080")`. 

The peer address should be a
resolvable address on your local-area-network `(e.g. "192.168.1.221:8080",
"localhost:8080", "Charles-MBP:8080")`

---

To talk to yourself, you can try:
* Listening Port: `8080`
* Peer Address: `localhost:8080`


# Additional Dependencies
* Requires a PortAudio (open-source audio I/O lib) installation

  **On MacOS:**
  ```
  brew install portaudio
  ```
 * Requires `pkg-config` a helper for compiling and linking installed libs.
 
   **On MacOS:**
   ```
   brew install pkg-config
   ```
