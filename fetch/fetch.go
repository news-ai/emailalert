package fetch

import (
	"fmt"
	"log"
	"runtime"

	"github.com/jprobinson/eazye"
	"github.com/news-ai/emailalert"
)

// https://github.com/golang/go/issues/3575 :(
var procs = runtime.NumCPU()

func FetchMail(cfg *emailalert.Config) {
	log.Print("getting mail")

	// give it 1000 buffer so we can load whatever IMAP throws at us in memory
	mail, err := eazye.GenerateUnread(cfg.MailboxInfo, cfg.MarkRead, false)
	if err != nil {
		log.Fatal("unable to get mail: ", err)
	}

	parseMessages(mail)
}

func parseMessages(mail chan eazye.Response) {
	fmt.Print("Hell")
	for resp := range mail {
		if resp.Err != nil {
			log.Fatalf("unable to fetch mail: %s", resp.Err)
			return
		}
		fmt.Print(resp.Email)
	}
}
