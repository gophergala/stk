package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Quota struct {
	QuotaMax       int `json:"quota_max"`
	QuotaRemaining int `json:"quota_remaining"`
}

type SearchResult struct {
	Items []struct {
		Tags  []string `json:"tags"`
		Owner struct {
			Reputation   int    `json:"reputation"`
			UserId       int    `json:"user_id"`
			UserType     string `json:"user_type"`
			ProfileImage string `json:"profile_image"`
			DisplayName  string `json:"display_name"`
			Link         string `json:"link"`
		} `json:"owner"`
		IsAnswered       bool   `json:"is_answered"`
		ViewCount        int    `json:"view_count"`
		AnswerCount      int    `json:"answer_count"`
		AcceptedAnswerId int    `json:"accepted_answer_id"`
		Score            int    `json:"score"`
		LastActivityDate int    `json:"last_activity_date"`
		CreationDate     int    `json:"creation_date"`
		QuestionId       int    `json:"question_id"`
		Link             string `json:"link"`
		Title            string `json:"title"`
	} `json:"items"`
	HasMore bool `json:"has_more"`
	Quota
}

type AnswerResponse struct {
	Items []struct {
		Owner struct {
			Reputation   int    `json:"reputation"`
			UserId       int    `json:"user_id"`
			UserType     string `json:"user_type"`
			AcceptRate   int    `json:"accept_rate"`
			ProfileImage string `json:"profile_image"`
			DisplayName  string `json:"display_name"`
			Link         string `json:"link"`
		} `json:"owner"`
		IsAccepted       bool   `json:"is_accepted"`
		Score            int    `json:"score"`
		LastActivityDate int    `json:"last_activity_date"`
		CreationDate     int    `json:"creation_date"`
		AnswerId         int    `json:"answer_id"`
		QuestionId       int    `json:"question_id"`
		Body             string `json:"body"`
	} `json:"items"`
	HasMore bool `json:"has_more"`
	Quota
}

func makeSearchRequest(query string) string {
	args := url.Values{}
	args.Set("order", "desc")
	args.Set("sort", "activity")
	args.Set("site", "stackoverflow")
	args.Set("accepted", "True")
	args.Set("q", query)

	return "https://api.stackexchange.com/2.2/search/advanced?" + args.Encode()
}

func makeAnswerRequest(answerIds ...int) string {
	var ids = make([]string, len(answerIds))

	for i, val := range answerIds {
		ids[i] = strconv.Itoa(val)
	}

	escaped := url.QueryEscape(strings.Join(ids, ";"))

	args := url.Values{}
	args.Set("order", "desc")
	args.Set("sort", "activity")
	args.Set("filter", "withbody")

	return "https://api.stackexchange.com/2.2/answers/" + escaped + "?" + args.Encode()
}

func get(url string) (body []byte, err error) {
	res, err := http.Get(url)
	defer res.Body.Close()

	if err != nil {
		return
	}

	return ioutil.ReadAll(res.Body)
}

func Search(query string) (result *SearchResult, err error) {
	url := makeSearchRequest(query)
	body, err := get(url)

	if err != nil {
		return nil, err
	}

	result = new(SearchResult)
	err = json.Unmarshal(body, &result)
	return
}

func GetAnswers(ids ...int) (result *AnswerResponse, err error) {
	url := makeAnswerRequest(ids...)
	body, err := get(url)

	if err != nil {
		return nil, err
	}

	result = new(AnswerResponse)
	err = json.Unmarshal(body, &result)
	return
}
