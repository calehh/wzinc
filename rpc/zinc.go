package rpc

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// query = `{
// 	"search_type": "match",
// 	"query":
// 	{
// 		"term": "DEMTSCHENKO",
// 		"start_time": "2021-06-02T14:28:31.894Z",
// 		"end_time": "2021-12-02T15:28:31.894Z"
// 	},
// 	"from": 0,
// 	"max_results": 20,
// 	"_source": []
// }`

type QueryReq struct {
	SearchType string `json:"search_type"`
	Query      Query  `json:"query"`
	From       int    `json:"from"`
	MaxResult  int    `json:"max_results"`
}

type Query struct {
	Term string `json:"term"`
}

var ErrQuery = errors.New("query err")

func (s *Service) zincQuery(query QueryReq, index string) ([]byte, error) {
	queryJson, _ := json.Marshal(&query)
	url := s.zincUrl + "/api/" + index + "/_search"
	req, err := http.NewRequest("POST", url, strings.NewReader(string(queryJson)))
	if err != nil {
		return nil, ErrQuery
	}
	req.SetBasicAuth(s.username, s.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, ErrQuery
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, ErrQuery
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrQuery
	}
	return body, nil
}
