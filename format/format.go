package format

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/news-ai/emailalert"
)

func FormatMail(cfg *emailalert.Config, sess *mgo.Session, t time.Time) {
	log.Print("formatting links")
	var results []emailalert.Tracking
	c := sess.DB("emailalert").C("keywordalerts")
	cKeyword := sess.DB("emailalert").C("keywordurl")

	err := c.Find(bson.M{"time": t}).All(&results)
	if err != nil {
		panic(err)
	}

	// row, column
	var csvData [][]string
	csvData = make([][]string, 0)

	if len(results) > 0 {
		var titleRepeated map[string]bool
		titleRepeated = make(map[string]bool)
		for _, result := range results {
			// Checks to see if title has already been added to CSV
			for _, href := range result.HREFs {
				title, url := getURLContent(href, cKeyword)
				if _, ok := titleRepeated[title]; !ok {
					log.Print(title)
					var columnData []string
					columnData = make([]string, 3)
					columnData[0] = result.Keyword
					columnData[1] = title
					columnData[2] = url
					csvData = append(csvData, columnData)
					titleRepeated[title] = true
				}
			}
		}
	}

	writeCSV(csvData)
}

func getURLContent(href string, cKeyword *mgo.Collection) (string, string) {
	splitHref := strings.Split(href, "http://")
	if len(splitHref) > 1 {
		href = "http://" + splitHref[1]
		splitHref = strings.Split(href, "&")
		href = splitHref[0]

		result := emailalert.Content{}
		err := cKeyword.Find(bson.M{"url": href}).One(&result)
		if err != nil {
			log.Print(err)
		}
		if result.Url == "" {
			resp, err := http.Get(href)
			if err != nil {
				log.Print(err)
			}
			root, err := html.Parse(resp.Body)
			if err != nil {
				log.Print(err)
			}

			title, ok := scrape.Find(root, scrape.ByTag(atom.Title))
			if ok {
				err = cKeyword.Insert(emailalert.Content{scrape.Text(title), href})
				return scrape.Text(title), href
			}
		}
	}
	return "", ""
}

func writeCSV(csvData [][]string) {
	file, err := os.Create("result.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)

	for _, value := range csvData {
		err := writer.Write(value)
		checkError("Cannot write to file", err)
	}

	defer writer.Flush()
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
