package investing

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type InvestingScheduleParser struct {
	idRegEx     *regexp.Regexp
	numberRegEx *regexp.Regexp
}

func NewInvestingScheduleParser() *InvestingScheduleParser {
	return &InvestingScheduleParser{
		idRegEx:     regexp.MustCompile(`(\d+)`),
		numberRegEx: regexp.MustCompile(`^-?\d+(\.\d+)?`),
	}
}

func (parser *InvestingScheduleParser) ParseScheduleHtml(s *goquery.Document, languageId int) (items []*InvestingScheduleRow, err error) {
	if s == nil {
		return nil, fmt.Errorf("argument html value is nil")
	}
	tableRows := s.Find("table tr[event_attr_id]")
	items = make([]*InvestingScheduleRow, len(tableRows.Nodes))
	tableRows.EachWithBreak(func(i int, s *goquery.Selection) bool {
		item, err := parser.parseScheduleRowHtml(s)
		if err != nil {
			return false
		}
		item.LanguageId = languageId
		items[i] = item
		return true
	})
	return
}

func (parser *InvestingScheduleParser) parseScheduleRowHtml(s *goquery.Selection) (*InvestingScheduleRow, error) {
	var err error
	result := InvestingScheduleRow{}

	if result.Id, err = parser.parseScheduleRowId(s); err != nil {
		return nil, err
	}
	if result.EventId, err = parseAttrValueToInt(s, "event_attr_id"); err != nil {
		return nil, err
	}
	if result.TimeStamp, err = parser.parseScheduleTimeStamp(s); err != nil {
		return nil, err
	}
	if result.Title, err = parser.parseScheduleTitle(s); err != nil {
		return nil, err
	}
	if result.CurrencyCode, err = parser.parseScheduleCurrencyCode(s); err != nil {
		return nil, err
	}
	if result.Sentiment, err = parser.parseScheduleSentiment(s); err != nil {
		return nil, err
	}
	if result.CountryName, err = parser.parseScheduleCountryName(s); err != nil {
		return nil, err
	}
	if result.Actual, err = parser.parseIndexValue(s, "act", "actual"); err != nil {
		return nil, err
	}
	if result.Forecast, err = parser.parseIndexValue(s, "fore", "forecast"); err != nil {
		return nil, err
	}
	if result.Previous, err = parser.parseIndexValue(s, "prev", "previous"); err != nil {
		return nil, err
	}
	if result.Type, err = parser.parseScheduleEventType(s); err != nil {
		return nil, err
	}

	return &result, err
}

func (parser *InvestingScheduleParser) parseScheduleRowId(s *goquery.Selection) (int, error) {
	idVal, err := getAttrValue(s, "id")
	if err != nil {
		return 0, nil
	}
	idStr := parser.idRegEx.FindString(idVal)
	if len(idStr) <= 0 {
		return 0, fmt.Errorf("id attribute has invalid value '%s'", idStr)
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (parser *InvestingScheduleParser) parseScheduleCountryName(s *goquery.Selection) (string, error) {
	flagCell := s.Find("span.ceFlags")
	return getAttrValue(flagCell, "title")
}

func (parser *InvestingScheduleParser) parseScheduleTimeStamp(s *goquery.Selection) (t time.Time, err error) {
	timeStr, err := getAttrValue(s, "data-event-datetime")
	if err != nil {
		return
	}

	t, err = time.Parse("2006/01/02 15:04:05", timeStr)
	if err != nil {
		return t, err
	}

	return
}

func (parser *InvestingScheduleParser) parseScheduleTitle(s *goquery.Selection) (string, error) {
	cell := s.Find("td.event a")
	if len(cell.Nodes) <= 0 {
		return "", fmt.Errorf("invalid html. title cell not found")
	}
	title := normalizeHtmlText(cell.Text())
	return title, nil
}

func (parser *InvestingScheduleParser) parseScheduleCurrencyCode(s *goquery.Selection) (code string, err error) {
	cell := s.Find("td.flagCur")
	if len(cell.Nodes) <= 0 {
		return "", fmt.Errorf("invalid html. currency cell not found")
	}
	code = normalizeHtmlText(cell.Text())
	return
}

func (parser *InvestingScheduleParser) parseScheduleSentiment(s *goquery.Selection) (sentiment int, err error) {
	items := s.Find("td.sentiment i.grayFullBullishIcon")
	sentiment = len(items.Nodes)
	if sentiment <= 0 || sentiment > 3 {
		return 0, fmt.Errorf("invalid html. sentiment has invalid value %d", sentiment)
	}
	return
}

func (parser *InvestingScheduleParser) parseIndexValue(s *goquery.Selection, className, fieldName string) (*float64, error) {
	cell := s.Find(fmt.Sprintf("td.%[1]s, td.%[1]s span", className))
	if len(cell.Nodes) <= 0 {
		return nil, fmt.Errorf("invalid html. %s cell not found", fieldName)
	}
	valueStr := normalizeHtmlText(cell.Text())
	if len(valueStr) <= 0 {
		return nil, nil
	}
	valueStr = parser.numberRegEx.FindString(valueStr)
	number, err := strconv.ParseFloat(valueStr, 64)

	if err != nil {
		return nil, err
	}

	return &number, err
}

func (parser *InvestingScheduleParser) parseScheduleEventType(s *goquery.Selection) (eventType ScheduleEventType, err error) {
	tag := s.Find("td.event span")
	if len(tag.Nodes) <= 0 {
		return Index, nil
	}
	typeStr, err := getAttrValue(tag, "data-img_key")
	if err != nil {
		return
	}
	switch typeStr {
	case "perliminary":
		return PreliminaryRelease, nil
	case "speach":
		return Speech, nil
	case "report":
		return Report, nil
	case "sandClock":
		return RetrievingData, nil
	}
	return Index, fmt.Errorf("invalid html. unknown event type %s", typeStr)
}
