package fetch

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
	"runtime"
	"strings"

	"golang.org/x/net/html"

	"github.com/jprobinson/eazye"
	"github.com/news-ai/emailalert"
)

// https://github.com/golang/go/issues/3575 :(
var procs = runtime.NumCPU()

var (
	anchorTag  = []byte("a")
	hrefAttr   = []byte("href")
	classAttr  = []byte("class")
	styleAttr  = []byte("style")
	httpPrefix = []byte("http")
	blank      = []byte("")

	urlRegex = regexp.MustCompile(`http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+`)
)

func FetchMail(cfg *emailalert.Config) {
	log.Print("getting mail")

	// give it 1000 buffer so we can load whatever IMAP throws at us in memory
	mail, err := eazye.GenerateAll(cfg.MailboxInfo, cfg.MarkRead, false)
	if err != nil {
		log.Fatal("unable to get mail: ", err)
	}

	parseMessages(mail)
}

func findHREFs(body []byte) []string {
	var hrefs []string

	z := html.NewTokenizer(bytes.NewReader(body))
loop:
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if err := z.Err(); err != nil && err != io.EOF {
				log.Print("unexpected error parsing html: ", err)
			}
			break loop
		case html.StartTagToken:
			tn, hasAttr := z.TagName()
			if bytes.Equal(tn, anchorTag) && hasAttr {
				// loop til we find an href attr or the end
				for {
					key, val, more := z.TagAttr()
					if bytes.Equal(hrefAttr, key) && bytes.HasPrefix(val, httpPrefix) {
						hrefs = append(hrefs, string(val))
						break
					}
					if !more {
						break
					}
				}
			}
		}
	}

	// found nothing? maybe regex for it?
	if len(hrefs) == 0 {
		matches := urlRegex.FindAll(body, -1)
		for _, match := range matches {
			hrefs = append(hrefs, string(match))
		}
	}
	return hrefs
}

func parseMessages(mail chan eazye.Response) {
	var keywordToEmails map[string][]eazye.Email
	keywordToEmails = make(map[string][]eazye.Email)
	for resp := range mail {
		if resp.Err != nil {
			log.Fatalf("unable to fetch mail: %s", resp.Err)
			return
		}
		// Grab keyword from the email subject
		keyword := strings.Replace(resp.Email.Subject, "Google Alert - ", "", -1)
		keyword = strings.Replace(keyword, "\"", "", -1)
		keywordToEmails[keyword] = append(keywordToEmails[keyword], resp.Email)
		// HTML := string(keywordToEmails[keyword][0].HTML[:])
		// fmt.Print(HTML)
		refs := findHREFs(resp.Email.HTML)
		fmt.Print(refs)
	}
}
