package tunnel

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

type rLoginClient struct {
	password string
	userName string
	termType string
}

// Read the RLogin client username, server username, and termtype parameters from the client
func rloginInit(localConn net.Conn, cfg config.Config) (*rLoginClient, error) {

	localConn.Write([]byte("\x00"))

	rlb := make([]byte, 512)
	_, err := localConn.Read(rlb)
	if err != nil {
		return nil, err
	}
	localConn.Write([]byte("\x00"))

	// Slice of "", client username, server username, terminal-type
	rloginData := strings.Split(string(rlb), "\x00")

	rlc := &rLoginClient{
		password: rloginData[1],
		termType: rloginData[3],
	}

	userNameRe := regexp.MustCompile(`^\[.+\].+`)
	// If the RLogin username contains a system tag, leave it as is
	if userNameRe.MatchString(rloginData[2]) {
		rlc.userName = rloginData[2]
		// Otherwise, prefix it with the system tag from the config file
	} else {
		// Strip any square brackets from the config file's system_tag value
		systemTagRe := regexp.MustCompile(`\[|\]`)
		systemTag := systemTagRe.ReplaceAllString(cfg.SystemTag, "")
		rlc.userName = fmt.Sprintf("[%s]%s", systemTag, rloginData[2])
	}

	return rlc, nil

}

// Connect to the SSH server and authenticate
func sshConnect(cfg config.Config) (*ssh.Client, error) {

	sshConfig := &ssh.ClientConfig{
		User: cfg.SSHUsername,
		Auth: []ssh.AuthMethod{ssh.Password(cfg.SSHPassword)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// Always accept key.
			return nil
		},
	}

	sshHost := fmt.Sprintf("%s:%s", cfg.SSHHost, cfg.SSHPort)
	serverConn, err := ssh.Dial("tcp", sshHost, sshConfig)

	return serverConn, err

}

// Connect to the RLogin server via the SSH connection
func startTunnel(cfg config.Config, rlc *rLoginClient, serverConn *ssh.Client) (net.Conn, error) {

	rloginHost := fmt.Sprintf("%s:%s", cfg.RLoginHost, cfg.RLoginPort)
	log.Printf("%s connecting to RLogin server %s via SSH tunnel", rlc.userName, rloginHost)

	remoteConn, err := serverConn.Dial("tcp", rloginHost)
	if err != nil {
		return remoteConn, err
	}

	// "Authenticate" and begin the actual RLogin session through the tunnel
	log.Printf("%s starting RLogin session, terminal type: %s", rlc.userName, rlc.termType)
	fmt.Fprintf(remoteConn, "\x00%s\x00%s\x00%s\x00", rlc.password, rlc.userName, rlc.termType)

	return remoteConn, err

}

func doTunnel(remoteConn net.Conn, localConn net.Conn) {
	go func() {
		io.Copy(remoteConn, localConn)
	}()
	io.Copy(localConn, remoteConn)
	localConn.Close()
}

// Start an SSH tunnel
func Start(localConn net.Conn, cfg config.Config) {

	defer localConn.Close()

	rlc, err := rloginInit(localConn, cfg)
	if err != nil {
		log.Printf("Error reading initial RLogin data from local client: %v", err)
		return
	}

	serverConn, err := sshConnect(cfg)
	if err != nil {
		log.Printf("%s encountered error connecting to SSH server: %v", rlc.userName, err)
		return
	}
	defer serverConn.Close()
	log.Printf("%s connected to SSH server %s:%s", rlc.userName, cfg.SSHHost, cfg.SSHPort)

	remoteConn, err := startTunnel(cfg, rlc, serverConn)
	if err != nil {
		log.Printf("%s encountered error connecting to RLogin server via SSH tunnel: %v", rlc.userName, err)
		return
	}
	defer remoteConn.Close()
	log.Printf("%s connected to RLogin server %s:%s via SSH tunnel", rlc.userName, cfg.RLoginHost, cfg.RLoginPort)

	doTunnel(remoteConn, localConn)

	log.Printf("%s disconnected", rlc.userName)

}
