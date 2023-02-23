package investing

import (
	"errors"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type InvestingCountryParser struct {
}

func (parser *InvestingCountryParser) parseCountryHtml(selection *goquery.Selection) (country *InvestingCountry, err error) {

	idValueStr, exists := selection.Find("input").Attr("value")

	if !exists {
		return nil, &ParsingError{
			Err: errors.New("invalid html missed value attribute or input tag"),
		}
	}

	result := InvestingCountry{}
	result.Id, err = strconv.Atoi(idValueStr)

	if err != nil {
		return nil, &ParsingError{Err: err}
	}

	result.Title = selection.Find("label").Text()

	if len(result.Title) <= 0 {
		return nil, &ParsingError{
			Err: errors.New("invalid html missed text or label tag"),
		}
	}

	return &result, nil
}

func (parser *InvestingCountryParser) ParseCountriesHtml(html *goquery.Document) (countries []*InvestingCountry, err error) {

	if html == nil {
		return nil, &ParsingError{
			Err: errors.New("argument html value is nil"),
		}
	}

	countriesHtml := html.Find("#filtersWrapper ul.countryOption li")

	if countriesHtml == nil || len(countriesHtml.Nodes) == 0 {
		return nil, &ParsingError{
			Err: errors.New("couldn't find country tags into html"),
		}
	}

	countries = make([]*InvestingCountry, len(countriesHtml.Nodes))

	countriesHtml.EachWithBreak(func(i int, s *goquery.Selection) bool {

		countries[i], err = parser.parseCountryHtml(s)

		return err == nil
	})

	return
}
