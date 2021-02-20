package server

import (
	"fmt"
	"log"
	"net"

	"github.com/echicken/dpc2/internal/config"
)

// Listen for connections
func Listen(cfg config.Config, handler func(net.Conn, config.Config)) {

	localInterface := fmt.Sprintf("%s:%s", cfg.LocalInterface, cfg.LocalPort)
	listener, err := net.Listen("tcp", localInterface)
	if err != nil {
		log.Fatalf("Bind error: %v", err)
		return
	}
	defer listener.Close()
	log.Printf("Listening on %s", localInterface)

	for {
		localConn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Listener error: %v", err)
			return
		}
		log.Print("Accepted connection")
		go handler(localConn, cfg)
	}

}
