package emailalert

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/jprobinson/eazye"
	"gopkg.in/mgo.v2"
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

func (c *Config) MgoSession() (*mgo.Session, error) {
	// make conn pass it to data
	sess, err := mgo.Dial("127.0.0.1")
	if err != nil {
		log.Printf("Unable to connect to emailalert db! - %s", err.Error())
		return sess, err
	}

	sess.SetMode(mgo.Eventual, true)
	return sess, nil
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
