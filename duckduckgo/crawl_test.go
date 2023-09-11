package duckduckgo_test

import (
	"testing"

	"github.com/dengchangdong/DuckDuckGo-API/duckduckgo"
	"github.com/dengchangdong/DuckDuckGo-API/typings"
)

func TestGet_html(t *testing.T) {
	search := typings.Search{
		Query:     "test",
		Region:    "",
		TimeRange: "",
	}
	response, err := duckduckgo.Get_results(search)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	t.Log(response)
}
