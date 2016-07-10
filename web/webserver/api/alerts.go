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
	fmt.Println("url: " + url)
	for i, href := range alert.HREFs {
		fmt.Println(href.Url)
		if href.Url == url {
			fmt.Println(href.Status)
			alert.HREFs[i].IsApproved = !alert.HREFs[i].IsApproved
			alert.HREFs[i].Status = !alert.HREFs[i].Status
		}
	}

	err := c.Update(bson.M{"_id": bson.ObjectIdHex(alert_id)}, bson.M{"$set": bson.M{"hrefs": alert.HREFs}})
	if err != nil {
		fmt.Printf("update fail %v\n", err)
	}

	return alert, nil
}

func SetAlertApprove(db *mgo.Database, alert_id string, url string, sentiment string) (emailalert.Gathering, error) {
	c := getNA(db)
	var alert emailalert.Gathering
	if err := c.FindId(bson.ObjectIdHex(alert_id)).One(&alert); err != nil {
		return alert, err
	}
	fmt.Println("url: " + url)
	for i, href := range alert.HREFs {
		fmt.Println(href.Url)
		if href.Url == url {
			fmt.Println(href.IsApproved)
			alert.HREFs[i].IsApproved = !alert.HREFs[i].IsApproved
			if sentiment == "positive" {
				alert.HREFs[i].NumPositive += 1
			} else if sentiment == "negative" {
				alert.HREFs[i].NumNegative += 1
			} else {
				alert.HREFs[i].NumNeutral += 1
			}
		}
	}

	err := c.Update(bson.M{"_id": bson.ObjectIdHex(alert_id)}, bson.M{"$set": bson.M{"hrefs": alert.HREFs}})
	if err != nil {
		fmt.Printf("update fail %v\n", err)
	}

	return alert, nil
}

func getNA(db *mgo.Database) *mgo.Collection {
	return db.C("gatheredalerts")
}
