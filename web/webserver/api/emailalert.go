// api package contains the Email Alert API.
package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jprobinson/go-utils/web"
	"gopkg.in/mgo.v2"

	"github.com/news-ai/emailalert"
)

var ErrDB = errors.New("problems accessing database")

// EmailAlertAPI is a struct that keeps a handle on the mgo session
type EmailAlertAPI struct {
	session *mgo.Session
}

// NewEmailAlertAPI creates a new EmailAlertAPI struct to run the emailalert API.
func NewEmailAlertAPI(config *emailalert.Config) *EmailAlertAPI {
	// make conn pass it to data
	sess, err := config.MgoSession()
	if err != nil {
		log.Fatalf("Unable to connect to emailalert db! - %s", err)
	}

	sess.SetMode(mgo.Eventual, true)
	return &EmailAlertAPI{sess}
}

func (n EmailAlertAPI) UrlPrefix() string {
	return "/v1"
}

func (n EmailAlertAPI) Handle(subRouter *mux.Router) {
	// ALERTS
	subRouter.HandleFunc("/get_all_articles", n.findAlerts).Methods("GET")
	subRouter.HandleFunc("/article_status/{alert_id}", n.setAlertStatus).Methods("GET")
	subRouter.HandleFunc("/article_approve/{alert_id}/{sentiment}", n.setAlertApprove).Methods("GET")
}

func (n EmailAlertAPI) findAlerts(w http.ResponseWriter, r *http.Request) {
	setCommonHeaders(w, r, "")

	s, db := n.getDB()
	defer s.Close()

	alerts, err := FindAlerts(db)
	if err != nil {
		log.Printf("Unable to access alerts! - %s", err.Error())
		web.ErrorResponse(w, ErrDB, http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, web.JsonResponseWrapper{Response: alerts})
}

func (n EmailAlertAPI) setAlertStatus(w http.ResponseWriter, r *http.Request) {
	setCommonHeaders(w, r, "")
	vars := mux.Vars(r)

	s, db := n.getDB()
	defer s.Close()

	alerts, err := SetAlertStatus(db, vars["alert_id"], string(r.URL.Query().Get("url")))
	if err != nil {
		log.Printf("Unable to access alerts! - %s", err.Error())
		web.ErrorResponse(w, ErrDB, http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, web.JsonResponseWrapper{Response: alerts})
}

func (n EmailAlertAPI) setAlertApprove(w http.ResponseWriter, r *http.Request) {
	setCommonHeaders(w, r, "")
	vars := mux.Vars(r)

	s, db := n.getDB()
	defer s.Close()

	alerts, err := SetAlertApprove(db, vars["alert_id"], string(r.URL.Query().Get("url")), vars["sentiment"])
	if err != nil {
		log.Printf("Unable to access alerts! - %s", err.Error())
		web.ErrorResponse(w, ErrDB, http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, web.JsonResponseWrapper{Response: alerts})
}

// setCommondHeaders is a utility function to set the 'Access-Control-Allow-Origin' to * and
// set the Content-Type to the given input. If not Content-Type is given, it defaults to
// 'application/json'.
func setCommonHeaders(w http.ResponseWriter, r *http.Request, contentType string) {
	origin := r.Header.Get("Origin")
	if len(origin) == 0 {
		origin = "*"
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if len(contentType) == 0 {
		w.Header().Set("Content-Type", web.JsonContentType)
	} else {
		w.Header().Set("Content-Type", contentType)
	}
}

func (n EmailAlertAPI) getDB() (*mgo.Session, *mgo.Database) {
	s := n.session.Copy()
	return s, s.DB("emailalert")
}
