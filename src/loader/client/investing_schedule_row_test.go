package client

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestParseScheduleRowId(t *testing.T) {
	tests := []struct {
		html   string
		result int
		err    error
	}{
		{
			html:   `<div id="row123456"></div>`,
			result: 123456,
			err:    nil,
		},
		{
			html:   `<div></div>`,
			result: 0,
			err:    &ParsingError{},
		},
		{
			html:   `<div id="rowABC"></div>`,
			result: 0,
			err:    &ParsingError{},
		},
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		selector := goquery.NewDocumentFromNode(node).Find("div")

		// Act
		value, err := parseScheduleRowId(selector)

		// Assert
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.result, value)
	}
}

func TestParseScheduleCountryName(t *testing.T) {
	tests := []struct {
		html   string
		result string
		err    error
	}{
		{
			html:   `<tr><td><span class="flagCur" title="China"></span></td></tr>`,
			result: "China",
			err:    nil,
		},
		{
			html:   `<tr><td><span class="flagCur"></span></td></tr>`,
			result: "",
			err:    &ParsingError{},
		},
		{
			html:   `<tr><td></td><td></td><td></td></tr>`,
			result: "",
			err:    &ParsingError{},
		},
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		selector := goquery.NewDocumentFromNode(node).Selection

		// Act
		value, err := parseScheduleCountryName(selector)

		// Assert
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.result, value)
	}
}
