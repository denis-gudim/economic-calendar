package investing

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func Test_InvestingCountryParser_ParseCountryHtml(t *testing.T) {
	tests := []struct {
		html           string
		expectedResult *InvestingCountry
		err            error
	}{
		{
			html: `<li><input value="123"><label>Text</label></li>`,
			expectedResult: &InvestingCountry{
				Id:         123,
				Title:      "Text",
				LanguageId: 0,
			},
			err: nil,
		},
		{
			html:           `<li><label>Text</label></li>`,
			expectedResult: nil,
			err:            &ParsingError{},
		},
		{
			html:           `<li><input><label>Text</label></li>`,
			expectedResult: nil,
			err:            &ParsingError{},
		},
		{
			html:           `<li><input value="abc"><label>Text</label></li>`,
			expectedResult: nil,
			err:            &ParsingError{},
		},
		{
			html:           `<li><input value="123"></li>`,
			expectedResult: nil,
			err:            &ParsingError{},
		},
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		html := goquery.NewDocumentFromNode(node)
		parser := &InvestingCountryParser{}

		// Act
		actualResult, err := parser.parseCountryHtml(html.Find("li"))

		// Assert
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expectedResult, actualResult)
	}
}

func Test_InvestingCountryParser_ParseCountriesHtml(t *testing.T) {

	tests := []struct {
		html           string
		expectedResult []*InvestingCountry
		err            error
	}{
		{
			html: `<ul class="countryOption">
						<li><input value="1"><label>Text 1</label></li>
						<li><input value="2"><label>Text 2</label></li>
					</ul>`,
			expectedResult: []*InvestingCountry{
				{Id: 1, Title: "Text 1", LanguageId: 0},
				{Id: 2, Title: "Text 2", LanguageId: 0},
			},
			err: nil,
		},
		{
			html:           `<div></div>`,
			expectedResult: nil,
			err:            &ParsingError{},
		},
		{
			html:           ``,
			expectedResult: nil,
			err:            &ParsingError{},
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

		parser := &InvestingCountryParser{}

		// Act
		actualResult, err := parser.ParseCountriesHtml(htmlDoc)

		// Assert
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expectedResult, actualResult)
	}
}
