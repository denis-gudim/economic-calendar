package investing

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestGetAttrValue(t *testing.T) {
	tests := []struct {
		html   string
		result string
		err    error
	}{
		{
			html:   `<div><input id="123"/></div>`,
			result: "123",
			err:    nil,
		},
		{
			html:   `<div><input /></div>`,
			result: "",
			err:    &ParsingError{},
		},
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		selector := goquery.NewDocumentFromNode(node).Find("input")

		// Act
		value, err := getAttrValue(selector, "id")

		// Assert
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.result, value)
	}
}

func TestParseAttrValueToInt(t *testing.T) {
	tests := []struct {
		html   string
		result int
		err    error
	}{
		{
			html:   `<div><input id="123"/></div>`,
			result: 123,
			err:    nil,
		},
		{
			html:   `<div><input/></div>`,
			result: 0,
			err:    &ParsingError{},
		},
		{
			html:   `<div><input id="abc"/></div>`,
			result: 0,
			err:    &ParsingError{},
		},
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		selector := goquery.NewDocumentFromNode(node).Find("input")

		// Act
		value, err := parseAttrValueToInt(selector, "id")

		// Assert
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.result, value)
	}
}

func Test_NormalizeHtmlText(t *testing.T) {
	tests := []struct {
		text           string
		expectedResult string
	}{
		{
			text:           `input&#039;s`,
			expectedResult: "input's",
		},
		{
			text:           `text&nbsp;`,
			expectedResult: "text",
		},
		{
			text:           `    text   `,
			expectedResult: "text",
		},
		{
			text:           ``,
			expectedResult: "",
		},
	}

	for _, test := range tests {
		// Arrange

		// Act
		actualResult := normalizeHtmlText(test.text)

		// Assert
		assert.Equal(t, test.expectedResult, actualResult)
	}
}
