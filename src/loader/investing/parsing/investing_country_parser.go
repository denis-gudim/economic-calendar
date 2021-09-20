package parsing

import (
	"economic-calendar/loader/investing/data"
	"errors"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type InvestingCountryParser struct {
}

func (parser *InvestingCountryParser) parseCountryHtml(selection *goquery.Selection) (*data.InvestingCountry, error) {

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

	return &data.InvestingCountry{Id: id, Title: title, LanguageId: 0}, nil
}

func (parser *InvestingCountryParser) ParseCountriesHtml(html *goquery.Document) (countries []*data.InvestingCountry, err error) {

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

	countries = make([]*data.InvestingCountry, len(countriesHtml.Nodes))

	countriesHtml.EachWithBreak(func(i int, s *goquery.Selection) bool {

		countries[i], err = parser.parseCountryHtml(s)

		return err == nil
	})

	return
}
