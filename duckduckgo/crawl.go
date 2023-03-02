package duckduckgo

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/acheong08/DuckDuckGo-API/types"
	"github.com/acheong08/DuckDuckGo-API/utils"
	"github.com/anaskhan96/soup"
)

func get_html(search types.Search) (string, error) {
	var base_url string = "html.duckduckgo.com"
	// POST form data
	var formdata = map[string]string{
		"q":  search.Query,
		"df": search.TimeRange,
		"kl": search.Region,
	}
	// URL encode form data
	var form string = utils.Url_encode(formdata)
	// Create POST request
	var request = http.Request{
		Method: "POST",
		URL: &url.URL{
			Host:   base_url,
			Path:   "/html/",
			Scheme: "https",
		},
		Header: map[string][]string{
			"Content-Type": {"application/x-www-form-urlencoded"},
			"Accept":       {"text/html"},
			"User-Agent":   {"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/111.0"},
		},
		Body: utils.StringToReadCloser(form),
	}
	// Send POST request
	var client = http.Client{}
	var response, err = client.Do(&request)
	if err != nil {
		return "", err
	}
	// Read response body
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	// Check status code
	if response.StatusCode != 200 {
		return "", errors.New("Status code: " + strconv.Itoa(response.StatusCode) + " Body: " + string(bodyBytes))
	}
	// Close response body
	err = response.Body.Close()
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}

func parse_html(html string) ([]types.Result, error) {
	// Results is an array of Result structs
	var final_results []types.Result = []types.Result{}
	// Parse
	doc := soup.HTMLParse(html)
	// Find each result__body
	result_bodies := doc.FindAll("div", "class", "result__body")
	// Loop through each result__body
	for _, item := range result_bodies {
		// Get text of result__title
		var title string = item.Find("a", "class", "result__a").FullText()
		// Get href of result__a
		var link string = item.Find("a", "class", "result__a").Attrs()["href"]
		// Get text of result__snippet
		var snippet string = item.Find("a", "class", "result__snippet").FullText()
		// Append to final_results
		final_results = append(final_results, types.Result{
			Title:   title,
			Link:    link,
			Snippet: snippet,
		})
	}
	return final_results, nil
}

func Get_results(search types.Search) ([]types.Result, error) {
	html, err := get_html(search)
	if err != nil {
		return nil, err
	}
	results, err := parse_html(html)
	if err != nil {
		return nil, err
	}
	return results, nil
}
