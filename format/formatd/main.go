package main

import (
	"flag"
	"log"

	"github.com/news-ai/emailalert"
	"github.com/news-ai/emailalert/format"

	"github.com/jprobinson/go-utils/utils"
	"gopkg.in/mgo.v2"
)

const logPath = "/var/log/emailalert/formatd.log"

var (
	logArg  = flag.String("log", logPath, "log path")
	reparse = flag.Bool("r", false, "reparse all alerts and events")
)

func main() {
	flag.Parse()

	if *logArg != "stderr" {
		logSetup := utils.NewDefaultLogSetup(*logArg)
		logSetup.SetupLogging()
		go utils.ListenForLogSignal(logSetup)
	} else {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	config := emailalert.NewConfig()
	sess, err := config.MgoSession()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	formatMail(config, sess)
}

func formatMail(config *emailalert.Config, sess *mgo.Session) {
	t := emailalert.GetTime()
	format.FormatMail(config, sess, t)
}
