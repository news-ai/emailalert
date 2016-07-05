package format

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/news-ai/emailalert"
)

func FormatMail(cfg *emailalert.Config, sess *mgo.Session, t time.Time) {
	log.Print("formatting links")
	var results []emailalert.Tracking
	c := sess.DB("emailalert").C("keywordalerts")

	err := c.Find(bson.M{"time": t}).All(&results)
	if err != nil {
		panic(err)
	}
	if len(results) > 0 {
		for _, result := range results {
			log.Print(result.Keyword)
		}
	}
}
