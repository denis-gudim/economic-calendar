package investing

import (
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

func (p *InvestingCalendarEventParser) ParseCalendarEventHtml(html *goquery.Document) (event *InvestingCalendarEvent, err error) {
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

	result := InvestingCalendarEvent{}

	result.Title, err = p.parseTitle(sectionTag)

	if err != nil {
		return
	}

	result.Overview, err = p.parseOverview(sectionTag)

	if err != nil {
		return
	}

	result.Source, result.SourceUrl, err = p.parseSourceInfo(sectionTag)

	if err != nil {
		return
	}

	result.Unit, err = p.parseUnit(sectionTag)

	if err != nil {
		return
	}

	result.Sentiment, err = p.parseSentiment(sectionTag)

	if err != nil {
		return
	}

	result.Country, err = p.parseCountry(sectionTag)

	return &result, err
}

func (p *InvestingCalendarEventParser) parseTitle(s *goquery.Selection) (title string, err error) {
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

func (p *InvestingCalendarEventParser) parseOverview(s *goquery.Selection) (overview string, err error) {
	tag := s.Find("#overViewBox div.left")

	if len(tag.Nodes) > 0 {
		overview = normalizeHtmlText(tag.Text())
	}

	return
}

func (p *InvestingCalendarEventParser) parseSourceInfo(s *goquery.Selection) (source string, sourceUrl string, err error) {
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

func (p *InvestingCalendarEventParser) parseUnit(s *goquery.Selection) (unit string, err error) {

	s.Find("#releaseInfo div.arial_14").EachWithBreak(func(i int, s *goquery.Selection) bool {

		indexValue := normalizeHtmlText(s.Text())

		if len(indexValue) <= 0 {
			return false
		}

		unit = p.unitRegEx.FindString(indexValue)

		return len(unit) > 0
	})

	return
}

func (p *InvestingCalendarEventParser) parseSentiment(s *goquery.Selection) (sentiment int, err error) {
	items := s.Find("i.grayFullBullishIcon")

	sentiment = len(items.Nodes)

	if sentiment <= 0 || sentiment > 3 {
		return 0, &ParsingError{
			Err: fmt.Errorf("invalid html. sentiment has invalid value %d", sentiment),
		}
	}

	return sentiment, nil
}

func (p *InvestingCalendarEventParser) parseCountry(s *goquery.Selection) (country string, err error) {
	tag := s.Find("i.ceFlags")

	if len(tag.Nodes) <= 0 {
		return "", &ParsingError{
			Err: errors.New("invalid html country tag not found"),
		}
	}

	country, err = getAttrValue(tag, "title")

	return
}
