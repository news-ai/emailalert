package api

import (
	"fmt"

	"github.com/news-ai/emailalert"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func FindAlerts(db *mgo.Database) ([]emailalert.Gathering, error) {
	c := getNA(db)
	var alerts []emailalert.Gathering
	if err := c.Find(bson.M{"time": emailalert.GetTime()}).All(&alerts); err != nil {
		return alerts, err
	}

	return alerts, nil
}

func SetAlertStatus(db *mgo.Database, alert_id string, url string) (emailalert.Gathering, error) {
	c := getNA(db)
	var alert emailalert.Gathering
	if err := c.FindId(bson.ObjectIdHex(alert_id)).One(&alert); err != nil {
		return alert, err
	}

	return alert, nil
}

func getNA(db *mgo.Database) *mgo.Collection {
	return db.C("gatheredalerts")
}
