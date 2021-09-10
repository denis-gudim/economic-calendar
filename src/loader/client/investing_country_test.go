package client

import (
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
	// Arrange
	htmlStr := `<html>
					<head><title>Title</title></head>
					<body>
						<div></div>
						<div id="filtersWrapper">
							<ul>
								<li><input value="1"><label>Text 1</label></li>
								<li><input value="2"><label>Text 2</label></li>
							</ul>
						</div>
					</body>
				</html>`
	htmlDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))

	// Act
	countres, err := parseCountriesHtml(htmlDoc)

	// Assert
	assert.Nil(t, err)
	assert.Len(t, countres, 2)

	assert.Equal(t, 1, countres[0].Id)
	assert.Equal(t, "Text 1", countres[0].Title)

	assert.Equal(t, 2, countres[1].Id)
	assert.Equal(t, "Text 2", countres[1].Title)
}
