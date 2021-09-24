package parsing

import (
	"economic-calendar/loader/investing/data"
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

func (parser *InvestingScheduleParser) ParseScheduleHtml(s *goquery.Document) (items []*data.InvestingScheduleRow, err error) {

	if s == nil {
		return nil, &ParsingError{
			Err: fmt.Errorf("argument html value is nil"),
		}
	}

	tableRows := s.Find("table tr[event_attr_id]")

	items = make([]*data.InvestingScheduleRow, len(tableRows.Nodes))

	tableRows.EachWithBreak(func(i int, s *goquery.Selection) bool {

		items[i], err = parser.parseScheduleRowHtml(s)

		return err == nil
	})

	return
}

func (parser *InvestingScheduleParser) parseScheduleRowHtml(s *goquery.Selection) (row *data.InvestingScheduleRow, err error) {

	result := data.InvestingScheduleRow{}

	result.Id, err = parser.parseScheduleRowId(s)

	if err != nil {
		return
	}

	result.EventId, err = parseAttrValueToInt(s, "event_attr_id")

	if err != nil {
		return
	}

	result.TimeStamp, err = parser.parseScheduleTimeStamp(s)

	if err != nil {
		return
	}

	result.Title, err = parser.parseScheduleTitle(s)

	if err != nil {
		return
	}

	result.CurrencyCode, err = parser.parseScheduleCurrencyCode(s)

	if err != nil {
		return
	}

	result.Sentiment, err = parser.parseScheduleSentiment(s)

	if err != nil {
		return
	}

	result.CountryName, err = parser.parseScheduleCountryName(s)

	if err != nil {
		return
	}

	result.Actual, err = parser.parseIndexValue(s, "act", "actual")

	if err != nil {
		return
	}

	result.Forecast, err = parser.parseIndexValue(s, "fore", "forecast")

	if err != nil {
		return
	}

	result.Previous, err = parser.parseIndexValue(s, "prev", "previous")

	if err != nil {
		return
	}

	result.Type, err = parser.parseScheduleEventType(s)

	return &result, err
}

func (parser *InvestingScheduleParser) parseScheduleRowId(s *goquery.Selection) (id int, err error) {

	idVal, err := getAttrValue(s, "id")

	if err != nil {
		return
	}

	idStr := parser.idRegEx.FindString(idVal)

	if len(idStr) <= 0 {
		return 0, &ParsingError{
			Err: fmt.Errorf("id attribute has invalid value '%s'", idStr),
		}
	}

	id, err = strconv.Atoi(idStr)

	if err != nil {
		return 0, &ParsingError{Err: err}
	}

	return
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
		return t, &ParsingError{Err: err}
	}

	return
}

func (parser *InvestingScheduleParser) parseScheduleTitle(s *goquery.Selection) (title string, err error) {
	cell := s.Find("td.event a")

	if len(cell.Nodes) <= 0 {
		return "", &ParsingError{
			Err: fmt.Errorf("invalid html. title cell not found"),
		}
	}

	title = normalizeHtmlText(cell.Text())

	return
}

func (parser *InvestingScheduleParser) parseScheduleCurrencyCode(s *goquery.Selection) (code string, err error) {
	cell := s.Find("td.flagCur")

	if len(cell.Nodes) <= 0 {
		return "", &ParsingError{
			Err: fmt.Errorf("invalid html. currency cell not found"),
		}
	}

	code = normalizeHtmlText(cell.Text())

	return
}

func (parser *InvestingScheduleParser) parseScheduleSentiment(s *goquery.Selection) (sentiment int, err error) {
	items := s.Find("td.sentiment i.grayFullBullishIcon")

	sentiment = len(items.Nodes)

	if sentiment <= 0 || sentiment > 3 {
		return 0, &ParsingError{
			Err: fmt.Errorf("invalid html. sentiment has invalid value %d", sentiment),
		}
	}

	return
}

func (parser *InvestingScheduleParser) parseIndexValue(s *goquery.Selection, className, fieldName string) (*float64, error) {
	cell := s.Find(fmt.Sprintf("td.%[1]s, td.%[1]s span", className))

	if len(cell.Nodes) <= 0 {
		return nil, &ParsingError{
			Err: fmt.Errorf("invalid html. %s cell not found", fieldName),
		}
	}

	valueStr := normalizeHtmlText(cell.Text())

	if len(valueStr) <= 0 {
		return nil, nil
	}

	valueStr = parser.numberRegEx.FindString(valueStr)

	number, err := strconv.ParseFloat(valueStr, 64)

	if err != nil {
		return nil, &ParsingError{Err: err}
	}

	return &number, err
}

func (parser *InvestingScheduleParser) parseScheduleEventType(s *goquery.Selection) (eventType data.ScheduleEventType, err error) {
	tag := s.Find("td.event span")

	if len(tag.Nodes) <= 0 {
		return data.Index, nil
	}

	typeStr, err := getAttrValue(tag, "data-img_key")

	if err != nil {
		return
	}

	switch typeStr {
	case "perliminary":
		return data.PreliminaryRelease, nil
	case "speach":
		return data.Speech, nil
	case "report":
		return data.Report, nil
	case "sandClock":
		return data.RetrievingData, nil
	}

	return data.Index, &ParsingError{
		Err: fmt.Errorf("invalid html. unknown event type %s", typeStr),
	}
}
