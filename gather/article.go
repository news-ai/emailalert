package gather

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type articleResp struct {
	Url               string   `json:"url"`
	Name              string   `json:"name"`
	Created_at        string   `json:"created_at"`
	Header_image      string   `json:"header_image"`
	Basic_summary     string   `json:"basic_summary"`
	Opening_paragraph string   `json:"opening_paragraph"`
	Keywords          []string `json:"keywords"`
	Authors           []string `json:"authors"`
}

func callArticleExtractor(url string) articleResp {
	var resp *http.Response
	var articleR articleResp

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
