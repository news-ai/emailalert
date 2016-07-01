package main

import (
	"flag"
	"log"
	"time"

	"github.com/jprobinson/go-utils/utils"

	"github.com/news-ai/emailalert"
	"github.com/news-ai/emailalert/fetch"
)

const logPath = "/var/log/emailalert/fetchd.log"

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

	fetchMail(config)
}

func fetchMail(config *emailalert.Config) {
	for {
		fetch.FetchMail(config)
		time.Sleep(30 * time.Second)
	}
}
