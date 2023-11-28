package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

const (
	GoogleSuggestUrl = "https://suggestqueries.google.com/complete/search?output=toolbar&hl=en&q="
)

// <?xml version="1.0"?><toplevel><CompleteSuggestion><suggestion data="testing"/></CompleteSuggestion><CompleteSuggestion><suggestion data="testing and commissioning"/></CompleteSuggestion><CompleteSuggestion><suggestion data="testing synonym"/></CompleteSuggestion><CompleteSuggestion><suggestion data="testing meaning"/></CompleteSuggestion><CompleteSuggestion><suggestion data="testing library"/></CompleteSuggestion><CompleteSuggestion><suggestion data="testing the waters"/></CompleteSuggestion><CompleteSuggestion><suggestion data="testing and certification"/></CompleteSuggestion><CompleteSuggestion><suggestion data="testing pyramid"/></CompleteSuggestion><CompleteSuggestion><suggestion data="testing the water meaning"/></CompleteSuggestion><CompleteSuggestion><suggestion data="test speed"/></CompleteSuggestion></toplevel>%

type Toplevel struct {
	XMLName            xml.Name `xml:"toplevel"`
	Text               string   `xml:",chardata"`
	CompleteSuggestion []struct {
		Text       string `xml:",chardata"`
		Suggestion struct {
			Text string `xml:",chardata"`
			Data string `xml:"data,attr"`
		} `xml:"suggestion"`
	} `xml:"CompleteSuggestion"`
}

func GoogleSuggest(query string) ([]list.Item, error) {
	items := []list.Item{}
	if query == "" {
		return items, nil
	}

	query = strings.ReplaceAll(query, " ", "+")
	url := GoogleSuggestUrl + query
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var suggestionsTopLevel Toplevel
	if err := xml.Unmarshal(b, &suggestionsTopLevel); err != nil {
		return nil, err
	}

	for _, s := range suggestionsTopLevel.CompleteSuggestion {
		items = append(items, item(s.Suggestion.Data))
	}

	return items, nil
}
