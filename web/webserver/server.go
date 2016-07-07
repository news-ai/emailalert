package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jprobinson/go-utils/utils"
	"github.com/jprobinson/go-utils/web"

	"github.com/news-ai/emailalert"
	"github.com/news-ai/emailalert/web/webserver/api"
)

func main() {
	config := emailalert.NewConfig()

	logSetup := utils.NewDefaultLogSetup(emailalert.ServerLog)
	logSetup.SetupLogging()
	go utils.ListenForLogSignal(logSetup)

	router := mux.NewRouter()

	api := api.NewEmailAlertAPI(config)
	apiRouter := router.PathPrefix(api.UrlPrefix()).Subrouter()
	api.Handle(apiRouter)

	staticRouter := router.PathPrefix("/").Subrouter()
	staticRouter.PathPrefix("/").Handler(http.FileServer(http.Dir(emailalert.WebDir)))

	handler := web.AccessLogHandler(emailalert.AccessLog, router)

	log.Fatal(http.ListenAndServe(":8080", handler))
}
