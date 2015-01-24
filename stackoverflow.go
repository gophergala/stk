package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	API_BASE = "https://api.stackexchange.com/2.2"
)

type CommonResponse struct {
	ErrorId      int    `json:"error_id"`
	ErrorMessage string `json:"error_message"`
	ErrorName    string `json:"error_name"`

	QuotaMax       int `json:"quota_max"`
	QuotaRemaining int `json:"quota_remaining"`

	HasMore bool `json:"has_more"`
}

type SearchResponse struct {
	CommonResponse
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
}

type AnswerResponse struct {
	CommonResponse
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
}

type ApiRequest struct {
	what   string
	ids    string
	params *url.Values
}

type Validator interface {
	IsValid() bool
	Error() error
}

func (res CommonResponse) IsValid() bool {
	if res.ErrorId > 0 || res.QuotaRemaining == 0 {
		return false
	}

	return true
}

func (err CommonResponse) Error() error {
	return fmt.Errorf("Error Id: %v, %v: %v", err.ErrorId, err.ErrorName, err.ErrorMessage)
}

func makeUrl(request *ApiRequest) string {
	var buf *bytes.Buffer
	buf = bytes.NewBufferString(API_BASE)
	buf.WriteByte('/')
	buf.WriteString(request.what)

	if len(request.ids) > 0 {
		buf.WriteByte('/')
		buf.WriteString(request.ids)
		buf.WriteByte('/')

	}

	buf.WriteByte('?')
	buf.WriteString(request.params.Encode())

	return buf.String()
}

func makeSearchRequest(query string) string {
	args := url.Values{}
	args.Set("order", "desc")
	args.Set("sort", "activity")
	args.Set("site", "stackoverflow")
	args.Set("accepted", "True")
	args.Set("q", query)

	return makeUrl(&ApiRequest{
		what:   "search/advanced",
		params: &args,
	})
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
	args.Set("site", "stackoverflow")
	args.Set("filter", "withbody")

	return makeUrl(&ApiRequest{
		what:   "answers",
		ids:    escaped,
		params: &args,
	})
}

func get(url string) (body []byte, err error) {
	res, err := http.Get(url)
	defer res.Body.Close()

	if err != nil {
		return
	}

	return ioutil.ReadAll(res.Body)
}

func load(url string, result Validator) (err error) {
	body, err := get(url)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, result)

	if !result.IsValid() {
		err = result.Error()
	}

	return
}

func Search(query string) (*SearchResponse, error) {
	url := makeSearchRequest(query)
	result := new(SearchResponse)
	err := load(url, result)

	return result, err
}

func GetAnswers(ids ...int) (*AnswerResponse, error) {
	url := makeAnswerRequest(ids...)
	result := new(AnswerResponse)
	err := load(url, result)

	return result, err
}
