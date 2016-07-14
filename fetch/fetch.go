package fetch

import (
	"bytes"
	"io"
	"log"
	"regexp"
	"runtime"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/news-ai/eazye"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

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

func FetchMail(cfg *emailalert.Config, sess *mgo.Session, t time.Time) {
	log.Print("getting mail")

	// give it 1000 buffer so we can load whatever IMAP throws at us in memory
	mail, err := eazye.GenerateUnreadSince(cfg.MailboxInfo, t, cfg.MarkRead, false)
	if err != nil {
		log.Fatal("unable to get mail: ", err)
	}

	parseMessages(mail, sess, t)
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func formatURL(href string) string {
	splitCharacter := "http://"
	if strings.Count(href, "https://") > 1 {
		splitCharacter = "https://"
	}
	splitHref := strings.Split(href, splitCharacter)
	splitHref = deleteEmpty(splitHref)
	if len(splitHref) > 1 {
		href = splitCharacter + splitHref[1]
		splitAmpersandHref := strings.Split(href, "&")
		return splitAmpersandHref[0]
	}
	return ""
}

func findHREFs(body []byte) []string {
	var hrefs []string
	var repeatHrefs map[string]bool
	repeatHrefs = make(map[string]bool)

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
						href := formatURL(string(val))
						if href != "" {
							if _, ok := repeatHrefs[href]; !ok {
								hrefs = append(hrefs, href)
								repeatHrefs[href] = true
							}
						}
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
			href := formatURL(string(match))
			if href != "" {
				hrefs = append(hrefs, href)
			}
		}
	}
	return hrefs
}

func parseMessages(mail chan eazye.Response, sess *mgo.Session, t time.Time) {
	var keywords map[string]bool
	keywords = make(map[string]bool)
	var keywordToEmails map[string][]eazye.Email
	keywordToEmails = make(map[string][]eazye.Email)
	var keywordToRefs map[string][]string
	keywordToRefs = make(map[string][]string)

	// Get information from email & extract links
	for resp := range mail {
		if resp.Err != nil {
			log.Fatalf("unable to fetch mail: %s", resp.Err)
			return
		}
		// Grab keyword from the email subject
		keyword := strings.Replace(resp.Email.Subject, "Google Alert - ", "", -1)
		keyword = strings.Replace(keyword, "\"", "", -1)
		keywords[keyword] = true
		log.Print("getting keyword: " + keyword)

		// Add email and HREFs to keywords
		keywordToEmails[keyword] = append(keywordToEmails[keyword], resp.Email)
		keywordToRefs[keyword] = findHREFs(resp.Email.HTML)
	}

	log.Println(t)
	// Add Keywords -> Links into MongoDB
	for keyword, _ := range keywords {
		track := emailalert.Tracking{keyword, keywordToRefs[keyword], t}
		log.Print("mongo keyword: " + keyword)
		c := sess.DB("emailalert").C("keywordalerts")

		result := emailalert.Tracking{}
		err := c.Find(bson.M{"keyword": keyword, "time": t}).One(&result)
		if err != nil {
			log.Print(err)
		}
		if result.Keyword == "" {
			err = c.Insert(&track)
			if err != nil {
				log.Print(err)
			}
			log.Print("Added Keyword " + keyword)
		} else {
			log.Print("Keyword " + keyword + " already exists")
		}
	}
}
