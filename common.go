package emailalert

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/jinzhu/now"
	"github.com/jprobinson/eazye"
	"gopkg.in/mgo.v2"
)

const (
	configFile = "/opt/emailalert/etc/config.json"

	ServerLog = "/var/log/emailalert/server.log"
	FetchLog  = "/var/log/emailalert/fetchd.log"
	AccessLog = "/var/log/emailalert/access.log"
)

type Content struct {
	Title    string
	Url      string
	Body     string
	Keywords []string
}

type Tracking struct {
	Keyword string
	HREFs   []string
	Time    time.Time
}

type Gathering struct {
	Keyword string
	HREFs   []Content
	Time    time.Time
}

type Sentence struct {
	Value   string   `json:"sentence"bson:"sentence"`
	Phrases []string `json:"noun_phrases"bson:"noun_phrases"`
}

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

func GetTime() time.Time {
	locname, offset := time.Now().Zone()
	loc := time.FixedZone(locname, offset)
	t := now.BeginningOfDay().In(loc)
	return t
}
