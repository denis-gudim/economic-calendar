package investing

import (
	"errors"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type InvestingCountryParser struct {
}

func (parser *InvestingCountryParser) parseCountryHtml(selection *goquery.Selection) (*InvestingCountry, error) {
	var err error
	idValueStr, exists := selection.Find("input").Attr("value")
	if !exists {
		return nil, errors.New("invalid html missed value attribute or input tag")
	}
	result := InvestingCountry{}
	if result.Id, err = strconv.Atoi(idValueStr); err != nil {
		return nil, err
	}
	if result.Title = selection.Find("label").Text(); len(result.Title) <= 0 {
		return nil, errors.New("invalid html missed text or label tag")
	}
	return &result, nil
}

func (parser *InvestingCountryParser) ParseCountriesHtml(html *goquery.Document) ([]*InvestingCountry, error) {
	if html == nil {
		return nil, errors.New("argument html value is nil")
	}
	countriesHtml := html.Find("#filtersWrapper ul.countryOption li")
	if countriesHtml == nil || len(countriesHtml.Nodes) == 0 {
		return nil, errors.New("couldn't find country tags into html")
	}
	countries := make([]*InvestingCountry, len(countriesHtml.Nodes))
	countriesHtml.EachWithBreak(func(i int, s *goquery.Selection) bool {
		country, err := parser.parseCountryHtml(s)
		if err != nil {
			return false
		}
		countries[i] = country
		return true
	})
	return countries, nil
}
