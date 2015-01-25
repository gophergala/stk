package stackoverflow

import (
	"bytes"
	"net/url"
	"strconv"
	"strings"
)

type SearchRequest struct {
	SiteID   string
	Sort     string
	Query    string
	Tags     []string
	Accepted bool
}

type AnswerRequest struct {
	AnswerIDS []int
	Sort      string
	SiteID    string
}

type apiRequest struct {
	what   string
	ids    string
	params *url.Values
}

func makeURL(request *apiRequest) string {
	var buf *bytes.Buffer
	buf = bytes.NewBufferString(APIBase)
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

func makeSearchRequest(request *SearchRequest) string {
	args := url.Values{}
	args.Set("order", "desc")
	args.Set("sort", request.Sort)
	args.Set("site", request.SiteID)
	args.Set("accepted", strconv.FormatBool(request.Accepted))
	args.Set("q", request.Query)

	return makeURL(&apiRequest{
		what:   "search/advanced",
		params: &args,
	})
}

func makeAnswerRequest(request *AnswerRequest) string {
	var ids = make([]string, len(request.AnswerIDS))

	for i, val := range request.AnswerIDS {
		ids[i] = strconv.Itoa(val)
	}

	escaped := url.QueryEscape(strings.Join(ids, ";"))

	args := url.Values{}
	args.Set("order", "desc")
	args.Set("filter", "withbody")
	args.Set("sort", request.Sort)
	args.Set("site", request.SiteID)

	return makeURL(&apiRequest{
		what:   "answers",
		ids:    escaped,
		params: &args,
	})
}
