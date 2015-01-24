package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	ADVANCED_SEARCH_URL = "https://api.stackexchange.com/2.2/search/advanced?"
)

type SearchResult struct {
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
	println(string(body))
}
