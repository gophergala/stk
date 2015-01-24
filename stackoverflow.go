package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	ADVANCED_SEARCH_URL = "https://api.stackexchange.com/2.2/search/advanced?"
)

type Item struct {
	link        string `json:"link"`
	title       string `json:"title"`
	question_id int64  `json:"int"`
}

type SearchResult struct {
	items           []Item `json:"items"`
	quota_max       int    `json:"quota_max"`
	quota_remaining int    `json:"quota_remaining"`
}

func Search(query string) {
	args := url.Values{}
	args.Set("order", "desc")
	args.Set("sort", "activity")
	args.Set("site", "stackoverflow")
	args.Set("q", query)

	url := ADVANCED_SEARCH_URL + args.Encode()
	log.Printf(url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var dat map[string]interface{}

	if err := json.Unmarshal(body, &dat); err != nil {
		panic(err)
	}

	for k, v := range dat {
		println(k)
		println(v)
	}

	//println(string(body))
}
