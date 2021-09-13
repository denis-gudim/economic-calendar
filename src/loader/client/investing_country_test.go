package client

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestParseCountryHtml(t *testing.T) {
	tests := []struct {
		html   string
		result *InvestingCountry
		err    error
	}{
		{
			html:   `<li><input value="123"><label>Text</label></li>`,
			result: &InvestingCountry{123, "Text", 0},
			err:    nil,
		},
		{
			html:   `<li><label>Text</label></li>`,
			result: nil,
			err:    &ParsingError{},
		},
		{
			html:   `<li><input><label>Text</label></li>`,
			result: nil,
			err:    &ParsingError{},
		},
		{
			html:   `<li><input value="abc"><label>Text</label></li>`,
			result: nil,
			err:    &ParsingError{},
		},
		{
			html:   `<li><input value="123"></li>`,
			result: nil,
			err:    &ParsingError{},
		},
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		html := goquery.NewDocumentFromNode(node)

		// Act
		country, err := parseCountryHtml(html.Find("li"))

		// Assert
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.result, country)
	}
}

func TestParseCountriesHtml(t *testing.T) {

	tests := []struct {
		html   string
		result []*InvestingCountry
		err    error
	}{
		{
			html: `<ul>
						<li><input value="1"><label>Text 1</label></li>
						<li><input value="2"><label>Text 2</label></li>
					</ul>`,
			result: []*InvestingCountry{
				{1, "Text 1", 0},
				{2, "Text 2", 0},
			},
			err: nil,
		},
		{
			html:   `<div></div>`,
			result: nil,
			err:    &ParsingError{},
		},
		{
			html:   ``,
			result: nil,
			err:    &ParsingError{},
		},
	}

	for _, test := range tests {

		// Arrange
		htmlStr := fmt.Sprintf(
			`<html>
					<head><title>Title</title></head>
					<body>
						<div></div>
						<div id="filtersWrapper">%s</div>
					</body>
				</html>`, test.html)

		var htmlDoc *goquery.Document

		if len(test.html) > 0 {
			htmlDoc, _ = goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
		} else {
			htmlDoc = nil
		}

		// Act
		countres, err := parseCountriesHtml(htmlDoc)

		// Assert
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.result, countres)
	}
}
