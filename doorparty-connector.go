package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"gopkg.in/ini.v1"
	"log"
	"net"
	"regexp"
	"strings"
)

// Ripped from StackOverflow or something
func connChan(conn net.Conn) chan []byte {
    c := make(chan []byte)
    go func() {
        b := make([]byte, 1024)
        for {
            n, err := conn.Read(b)
            if n > 0 {
                res := make([]byte, n)
                // Copy the buffer so it doesn't get changed while read by the recipient.
                copy(res, b[:n])
                c <- res
            }
            if err != nil {
                c <- nil
                break
            }
        }
    }()
    return c
}

// Returns magic booleans so we'll know which side errored / disconnected
func pipeConns(c1 net.Conn, c2 net.Conn) bool {
	ch1 := connChan(c1)
	ch2 := connChan(c2)
	for {
		select {
			case b1 := <- ch1:
				if b1 == nil {
					return false
				} else {
					c2.Write(b1)
				}
			case b2 := <- ch2:
				if b2 == nil {
					return true
				} else {
					c1.Write(b2)
				}
		}
	}
}

func doTunnel(localConn net.Conn, cfg *ini.File) {
	
	defer localConn.Close()

	// Strip any square brackets from the config file's system_tag value
	systemTagRe := regexp.MustCompile(`\[|\]`)
	systemTag := systemTagRe.ReplaceAllString(cfg.Section("").Key("system_tag").Value(), "")

	// Read the RLogin client username, server username, and termtype parameters from the client
	rloginInit := make([]byte, 512)
	_, err := localConn.Read(rloginInit)
	if err != nil {
		log.Printf("Error reading initial RLogin data from local client")
		return
	}
	// Slice of "", client username, server username, terminal-type
	rloginData := strings.Split(string(rloginInit), "\x00")

	// Strip the system tag from the username in case the BBS included it
	userNameRe := regexp.MustCompile(`^\[.*\]`)
	userName := userNameRe.ReplaceAllString(rloginData[2], "")

	// Connect to the SSH server and authenticate
	sshConfig := &ssh.ClientConfig{
		User: cfg.Section("").Key("ssh_username").Value(),
		Auth: []ssh.AuthMethod{ssh.Password(cfg.Section("").Key("ssh_password").Value())},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// Always accept key.
			return nil
		},
	}
	sshHost := fmt.Sprintf("%s:%s", cfg.Section("").Key("ssh_host").Value(), cfg.Section("").Key("ssh_port").Value())
	log.Printf("%s connecting to SSH server %s", userName, sshHost)
	serverConn, err := ssh.Dial("tcp", sshHost, sshConfig)
	if err != nil {
		log.Printf("%s encountered error connecting to SSH server: %v", userName, err)
		return
	}
	defer serverConn.Close()
	log.Printf("%s connected to SSH server %s", userName, sshHost)

	// Connect to the RLogin server via the SSH connection
	rloginHost := fmt.Sprintf("%s:%s", cfg.Section("").Key("rlogin_host").Value(), cfg.Section("").Key("rlogin_port").Value())
	log.Printf("%s connecting to RLogin server %s via SSH tunnel", userName, rloginHost)
	remoteConn, err := serverConn.Dial("tcp", rloginHost)
	if err != nil {
		log.Printf("%s encountered error connecting to RLogin server via SSH tunnel: %v", userName, err)
		return
	}
	defer remoteConn.Close()
	log.Printf("%s connected to RLogin server %s via SSH tunnel", userName, rloginHost);

	// "Authenticate" and begin the actual RLogin session through the tunnel
	log.Printf("%s starting RLogin session for [%s]%s, terminal type: %s", userName, systemTag, userName, rloginData[3])
	fmt.Fprintf(remoteConn, "\x00%s\x00[%s]%s\x00%s\x00", rloginData[1], systemTag, userName, rloginData[3])

	// Pipe data between local client and remote RLogin server until one of them closes/errors
	remoteClosed := pipeConns(localConn, remoteConn)
	if (remoteClosed) {
		log.Printf("%s disconnected by remote server", userName)
	} else {
		log.Printf("%s disconnected locally", userName)
	}

}

func init() {
	log.SetFlags(log.Ldate|log.Ltime)
	log.Print("Initialized")
}

func main() {

	cfg, err := ini.Load("doorparty-connector.ini")
	if err != nil {
		log.Fatalf("Error reading doorparty-connector.ini: %v", err)
	}

	localInterface := fmt.Sprintf("%s:%s", cfg.Section("").Key("local_interface").Value(), cfg.Section("").Key("local_port").Value())
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