package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"strings"

	"github.com/echicken/dpc2/internal/config"
	"golang.org/x/crypto/ssh"
)

func doTunnel(localConn net.Conn, cfg config.Config) {

	defer localConn.Close()

	localConn.Write([]byte("\x00"))

	// Read the RLogin client username, server username, and termtype parameters from the client
	rloginInit := make([]byte, 512)
	_, err := localConn.Read(rloginInit)
	if err != nil {
		log.Printf("Error reading initial RLogin data from local client: %v %v", rloginInit, err)
		return
	}
	localConn.Write([]byte("\x00"))

	// Slice of "", client username, server username, terminal-type
	rloginData := strings.Split(string(rloginInit), "\x00")

	userName := ""
	userNameRe := regexp.MustCompile(`^\[.+\].+`)
	// If the RLogin username contains a system tag, leave it as is
	if userNameRe.MatchString(rloginData[2]) {
		userName = rloginData[2]
		// Otherwise, prefix it with the system tag from the config file
	} else {
		// Strip any square brackets from the config file's system_tag value
		systemTagRe := regexp.MustCompile(`\[|\]`)
		systemTag := systemTagRe.ReplaceAllString(cfg.SystemTag, "")
		userName = fmt.Sprintf("[%s]%s", systemTag, rloginData[2])
	}

	// Connect to the SSH server and authenticate
	sshConfig := &ssh.ClientConfig{
		User: cfg.SSHUsername,
		Auth: []ssh.AuthMethod{ssh.Password(cfg.SSHPassword)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// Always accept key.
			return nil
		},
	}
	sshHost := fmt.Sprintf("%s:%s", cfg.SSHHost, cfg.SSHPort)
	log.Printf("%s connecting to SSH server %s", userName, sshHost)
	serverConn, err := ssh.Dial("tcp", sshHost, sshConfig)
	if err != nil {
		log.Printf("%s encountered error connecting to SSH server: %v", userName, err)
		return
	}
	defer serverConn.Close()
	log.Printf("%s connected to SSH server %s", userName, sshHost)

	// Connect to the RLogin server via the SSH connection
	rloginHost := fmt.Sprintf("%s:%s", cfg.RLoginHost, cfg.RLoginPort)
	log.Printf("%s connecting to RLogin server %s via SSH tunnel", userName, rloginHost)
	remoteConn, err := serverConn.Dial("tcp", rloginHost)
	if err != nil {
		log.Printf("%s encountered error connecting to RLogin server via SSH tunnel: %v", userName, err)
		return
	}
	defer remoteConn.Close()
	log.Printf("%s connected to RLogin server %s via SSH tunnel", userName, rloginHost)

	// "Authenticate" and begin the actual RLogin session through the tunnel
	log.Printf("%s starting RLogin session, terminal type: %s", userName, rloginData[3])
	fmt.Fprintf(remoteConn, "\x00%s\x00%s\x00%s\x00", rloginData[1], userName, rloginData[3])

	go func() {
		io.Copy(remoteConn, localConn)
	}()
	io.Copy(localConn, remoteConn)

	localConn.Close()
	log.Printf("%s disconnected", userName)

}

func init() {
	log.SetFlags(log.Ldate | log.Ltime)
	log.Print("Initialized")
}

func main() {

	cfg := config.Get()

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
		go doTunnel(localConn, cfg)
	}

}
