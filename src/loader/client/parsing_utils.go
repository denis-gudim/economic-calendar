package client

import (
	"fmt"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func getAttrValue(s *goquery.Selection, attrName string) (value string, err error) {
	if s == nil {
		return "", &ParsingError{
			Err: fmt.Errorf("argument html value is nil"),
		}
	}

	value, exists := s.Attr(attrName)

	if !exists {
		return "", &ParsingError{
			Err: fmt.Errorf("html invalid. '%s' attribute is missed", attrName),
		}
	}

	return
}

func parseAttrValueToInt(s *goquery.Selection, attrName string) (value int, err error) {

	valueStr, err := getAttrValue(s, attrName)

	if err != nil {
		return
	}

	value, err = strconv.Atoi(valueStr)

	if err != nil {
		return 0, &ParsingError{Err: err}
	}

	return
}
