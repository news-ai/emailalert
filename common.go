package emailalert

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/jinzhu/now"
	"github.com/news-ai/eazye"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	configFile = "/opt/emailalert/etc/config.json"

	ServerLog = "/var/log/emailalert/server.log"
	AccessLog = "/var/log/emailalert/access.log"

	WebDir = "/opt/newshound/www"
)

type Content struct {
	Url               string     `json:"url"`
	Name              string     `json:"name"`
	Created_at        string     `json:"created_at"`
	Header_image      string     `json:"header_image"`
	Basic_summary     string     `json:"basic_summary"`
	Opening_paragraph string     `json:"opening_paragraph"`
	Keywords          []string   `json:"keywords"`
	Authors           []string   `json:"authors"`
	Html              string     `json:"-"`
	Text              string     `json:"text"`
	Tags              []string   `json:"tags"`
	Sentences         []Sentence `json:"sentences"`
	TopSentence       string     `json:"topsentence"`
	Status            bool       `json:"status"`
	IsApproved        bool       `json:"approved"`
	NumPositive       int        `json:"numpositive"`
	NumNegative       int        `json:"numnegative"`
	NumNeutral        int        `json:"numneutral"`
}

type Tracking struct {
	Keyword string
	HREFs   []string
	Time    time.Time
}

type Gathering struct {
	Id      bson.ObjectId "_id,omitempty"
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
