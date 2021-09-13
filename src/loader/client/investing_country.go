package client

import (
	"errors"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type InvestingCountry struct {
	Id         int
	Title      string
	LanguageId int
}

func parseCountryHtml(selection *goquery.Selection) (*InvestingCountry, error) {

	idValueStr, exists := selection.Find("input").Attr("value")

	if !exists {
		return nil, &ParsingError{
			Err: errors.New("invalid html missed value attribute or input tag"),
		}
	}

	id, err := strconv.Atoi(idValueStr)

	if err != nil {
		return nil, &ParsingError{Err: err}
	}

	title := selection.Find("label").Text()

	if len(title) <= 0 {
		return nil, &ParsingError{
			Err: errors.New("invalid html missed text or label tag"),
		}
	}

	return &InvestingCountry{id, title, 0}, nil
}

func parseCountriesHtml(html *goquery.Document) (countries []*InvestingCountry, err error) {

	if html == nil {
		return nil, &ParsingError{
			Err: errors.New("argument html value is nil"),
		}
	}

	countriesHtml := html.Find("#filtersWrapper ul li")

	if countriesHtml == nil || len(countriesHtml.Nodes) == 0 {
		return nil, &ParsingError{
			Err: errors.New("couldn't find country tags into html"),
		}
	}

	countries = make([]*InvestingCountry, len(countriesHtml.Nodes))

	countriesHtml.EachWithBreak(func(i int, s *goquery.Selection) bool {

		countries[i], err = parseCountryHtml(s)

		return err == nil
	})

	return
}
