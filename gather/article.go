package gather

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/news-ai/emailalert"
)

func callArticleExtractor(url string) emailalert.Content {
	var resp *http.Response
	var articleR emailalert.Content

	resp, err := http.Post("http://127.0.0.1:1030/", "application/json", bytes.NewReader([]byte(url)))
	if err != nil {
		log.Print("unable to hit article_extractor: ", err)
		return articleR
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&articleR); err != nil {
		return articleR
	}
	return articleR
}
