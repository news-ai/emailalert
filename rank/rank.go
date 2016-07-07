package rank

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/news-ai/emailalert"
)

func RankAlerts(cfg *emailalert.Config, sess *mgo.Session, t time.Time) {
	log.Print("ranking links: " + t.String())
	var results []emailalert.Gathering
	alertSession := sess.DB("emailalert")
	gatheredCollection := alertSession.C("gatheredalerts")
	// rankedCollection := alertSession.C("rankalerts")

	err := gatheredCollection.Find(bson.M{"time": t}).All(&results)
	if err != nil {
		log.Println(err)
	}
	if len(results) > 0 {
		for _, result := range results {
			log.Println(result.Keyword)
			for _, href := range result.HREFs {
				log.Print(href.Name + " " + href.Url)
			}
		}
	}
}
