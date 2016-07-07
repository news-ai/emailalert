package format

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/news-ai/emailalert"
)

func FormatMail(cfg *emailalert.Config, sess *mgo.Session, t time.Time) {
	log.Print("formatting links")
	var results []emailalert.Gathering
	c := sess.DB("emailalert").C("rankalerts")

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
			// Checks to see if content.Name has already been added to CSV
			for _, content := range result.HREFs {
				if _, ok := titleRepeated[content.Name]; !ok {
					log.Print(content.Name)
					var columnData []string
					columnData = make([]string, 6)
					columnData[0] = result.Keyword
					columnData[1] = content.Name
					columnData[2] = content.Url
					columnData[3] = content.Basic_summary
					columnData[4] = strings.Join(content.Keywords, ", ")
					columnData[5] = strings.Join(content.Tags, ", ")
					csvData = append(csvData, columnData)
					titleRepeated[content.Name] = true
				}
			}
		}
	}

	writeCSV(csvData)
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
