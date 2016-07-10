package gather

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/news-ai/emailalert"
)

func GatherAlerts(cfg *emailalert.Config, sess *mgo.Session, t time.Time) {
	log.Print("gathering links: " + t.String())
	var results []emailalert.Tracking
	alertSession := sess.DB("emailalert")
	alertsCollection := alertSession.C("keywordalerts")
	gatheredCollection := alertSession.C("gatheredalerts")

	err := alertsCollection.Find(bson.M{"time": t}).All(&results)
	if err != nil {
		log.Println(err)
	}
	if len(results) > 0 {
		for _, result := range results {
			log.Println(result.Keyword)

			singleKeyword := emailalert.Gathering{}
			singleKeyword.Keyword = result.Keyword
			singleKeyword.Time = t
			singleKeyword.HREFs = []emailalert.Content{}

			for _, href := range result.HREFs {
				article := callArticleExtractor(href)
				article.Tags, article.Sentences, article.TopSentence, err = callNP([]byte(article.Text))
				if err != nil {
					log.Println(err)
				}
				if article.Url == "" {
					article.Url = href
				}
				singleKeyword.HREFs = append(singleKeyword.HREFs, article)
			}
			result := emailalert.Gathering{}
			err := gatheredCollection.Find(bson.M{"keyword": result.Keyword, "time": t}).One(&result)
			if err != nil {
				log.Print(err)
			}
			if result.Keyword == "" {
				err = gatheredCollection.Insert(&singleKeyword)
				if err != nil {
					log.Print(err)
				}
				log.Print("Added Keyword " + result.Keyword)
			} else {
				log.Print("Keyword " + result.Keyword + " already exists")
			}
		}
	}
}
