package config

import (
	"log"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

// Config ...
type Config struct {
	SystemTag      string
	SSHUsername    string
	SSHPassword    string
	LocalInterface string
	LocalPort      string
	SSHHost        string
	SSHPort        string
	RLoginHost     string
	RLoginPort     string
}

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

// Get a Config
func Get() Config {

	var fn string
	if len(os.Args) >= 2 {
		fn = os.Args[1]
	} else {
		fn = "doorparty-connector.ini"
	}

	iniFile, err := ini.LooseLoad(fn)
	if err != nil {
		log.Fatalf("Error reading %v: %v", fn, err)
	}

	var cfg = Config{
		SystemTag:      getSetting("SYSTEM_TAG", iniFile, ""),
		SSHUsername:    getSetting("SSH_USERNAME", iniFile, ""),
		SSHPassword:    getSetting("SSH_PASSWORD", iniFile, ""),
		LocalInterface: getSetting("LOCAL_INTERFACE", iniFile, "0.0.0.0"),
		LocalPort:      getSetting("LOCAL_PORT", iniFile, "9999"),
		SSHHost:        getSetting("SSH_HOST", iniFile, "dp.throwbackbbs.com"),
		SSHPort:        getSetting("SSH_PORT", iniFile, "2022"),
		RLoginHost:     getSetting("RLOGIN_HOST", iniFile, "dp.throwbackbbs.com"),
		RLoginPort:     getSetting("RLOGIN_PORT", iniFile, "513"),
	}

	return cfg

}
