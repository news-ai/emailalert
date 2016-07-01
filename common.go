package emailalert

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/jprobinson/eazye"
)

const (
	configFile = "/opt/emailalert/etc/config.json"

	ServerLog = "/var/log/emailalert/server.log"
	FetchLog  = "/var/log/emailalert/fetchd.log"
	AccessLog = "/var/log/emailalert/access.log"
)

type Config struct {
	MarkRead          bool `json:"mark_as_read"`
	eazye.MailboxInfo `,inline`
}

func NewConfig() *Config {
	config := Config{}

	readBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Cannot read config file: %s %s", config, err)
	}

	err = json.Unmarshal(readBytes, &config)
	if err != nil {
		log.Fatalf("Cannot parse JSON in config file: %s %s", config, err)
	}

	return &config
}
