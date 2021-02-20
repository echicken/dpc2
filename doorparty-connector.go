package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh"
	"gopkg.in/ini.v1"
)

func getSetting(val string, iniFile *ini.File, defaultValue string) string {
	var ret = os.Getenv(val)
	if ret == "" {
		ret = iniFile.Section("").Key(strings.ToLower(val)).Value()
	}
	if ret == "" {
		return defaultValue
	}
	return ret
}

type config struct {
	systemTag      string
	sshUsername    string
	sshPassword    string
	localInterface string
	localPort      string
	sshHost        string
	sshPort        string
	rloginHost     string
	rloginPort     string
}

func getConfig(iniFile *ini.File) config {
	var cfg = config{}
	cfg.systemTag = getSetting("SYSTEM_TAG", iniFile, "")
	cfg.sshUsername = getSetting("SSH_USERNAME", iniFile, "")
	cfg.sshPassword = getSetting("SSH_PASSWORD", iniFile, "")
	cfg.localInterface = getSetting("LOCAL_INTERFACE", iniFile, "0.0.0.0")
	cfg.localPort = getSetting("LOCAL_PORT", iniFile, "513")
	cfg.sshHost = getSetting("SSH_HOST", iniFile, "dp.throwbackbbs.com")
	cfg.sshPort = getSetting("SSH_PORT", iniFile, "2022")
	cfg.rloginHost = getSetting("RLOGIN_HOST", iniFile, "dp.throwbackbbs.com")
	cfg.rloginPort = getSetting("RLOGIN_PORT", iniFile, "513")
	return cfg
}

func doTunnel(localConn net.Conn, cfg config) {

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
		systemTag := systemTagRe.ReplaceAllString(cfg.systemTag, "")
		userName = fmt.Sprintf("[%s]%s", systemTag, rloginData[2])
	}

	// Connect to the SSH server and authenticate
	sshConfig := &ssh.ClientConfig{
		User: cfg.sshUsername,
		Auth: []ssh.AuthMethod{ssh.Password(cfg.sshPassword)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// Always accept key.
			return nil
		},
	}
	sshHost := fmt.Sprintf("%s:%s", cfg.sshHost, cfg.sshPort)
	log.Printf("%s connecting to SSH server %s", userName, sshHost)
	serverConn, err := ssh.Dial("tcp", sshHost, sshConfig)
	if err != nil {
		log.Printf("%s encountered error connecting to SSH server: %v", userName, err)
		return
	}
	defer serverConn.Close()
	log.Printf("%s connected to SSH server %s", userName, sshHost)

	// Connect to the RLogin server via the SSH connection
	rloginHost := fmt.Sprintf("%s:%s", cfg.rloginHost, cfg.rloginPort)
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

	var fn string
	if len(os.Args) >= 2 {
		fn = os.Args[1]
	} else {
		fn = "doorparty-connector.ini"
	}

	sf, err := ini.LooseLoad(fn)
	if err != nil {
		log.Fatalf("Error reading %v: %v", fn, err)
	}

	cfg := getConfig(sf)

	localInterface := fmt.Sprintf("%s:%s", cfg.localInterface, cfg.localPort)
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
