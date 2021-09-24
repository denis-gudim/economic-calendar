package parsing

import (
	"economic-calendar/loader/investing/data"
	"errors"
	"fmt"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

type InvestingCalendarEventParser struct {
	unitRegEx *regexp.Regexp
}

func NewInvestingCalendarEventParser() *InvestingCalendarEventParser {
	return &InvestingCalendarEventParser{
		unitRegEx: regexp.MustCompile(`[^(-?\d+(\.\d)?)]`),
	}
}

func (parser *InvestingCalendarEventParser) ParseCalendarEventHtml(html *goquery.Document) (event *data.InvestingCalendarEvent, err error) {
	if html == nil {
		return nil, &ParsingError{
			Err: fmt.Errorf("argument html value is nil"),
		}
	}

	sectionTag := html.Find("#leftColumn")

	if len(sectionTag.Nodes) <= 0 {
		return nil, &ParsingError{
			Err: fmt.Errorf("invalid html. couldn't find section details node"),
		}
	}

	result := data.InvestingCalendarEvent{}

	result.Title, err = parser.parseTitle(sectionTag)

	if err != nil {
		return
	}

	result.Overview, err = parser.parseOverview(sectionTag)

	if err != nil {
		return
	}

	result.Source, result.SourceUrl, err = parser.parseSourceInfo(sectionTag)

	if err != nil {
		return
	}

	result.Unit, err = parser.parseUnit(sectionTag)

	return &result, err
}

func (parser *InvestingCalendarEventParser) parseTitle(s *goquery.Selection) (title string, err error) {
	tag := s.Find("h1.ecTitle")

	if len(tag.Nodes) <= 0 {
		return "", &ParsingError{
			Err: errors.New("invalid html missed title header tag"),
		}
	}

	title = normalizeHtmlText(tag.Text())

	if len(title) <= 0 {
		return "", &ParsingError{
			Err: errors.New("invalid html title string is empty"),
		}
	}

	return
}

func (parser *InvestingCalendarEventParser) parseOverview(s *goquery.Selection) (overview string, err error) {
	tag := s.Find("#overViewBox div.left")

	if len(tag.Nodes) > 0 {
		overview = normalizeHtmlText(tag.Text())
	}

	return
}

func (parser *InvestingCalendarEventParser) parseSourceInfo(s *goquery.Selection) (source string, sourceUrl string, err error) {
	tag := s.Find("div.right div:last-child a")

	if len(tag.Nodes) <= 0 {
		return
	}

	source, err = getAttrValue(tag, "title")

	if err != nil {
		return
	}

	sourceUrl, err = getAttrValue(tag, "href")

	return
}

func (parser *InvestingCalendarEventParser) parseUnit(s *goquery.Selection) (unit string, err error) {

	s.Find("#releaseInfo div.arial_14").EachWithBreak(func(i int, s *goquery.Selection) bool {

		indexValue := normalizeHtmlText(s.Text())

		if len(indexValue) <= 0 {
			return false
		}

		unit = parser.unitRegEx.FindString(indexValue)

		return len(unit) > 0
	})

	return
}
