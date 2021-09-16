package client

import (
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestParseScheduleRowHtml(t *testing.T) {
	tests := []struct {
		html   string
		result *InvestingScheduleRow
		err    error
	}{
		{
			html: `<table>
						<tr id="eventRowId_436932" class="js-event-item" event_attr_id="377" data-event-datetime="2021/09/16 08:00:00">
							<td class="first left time js-time" title="Event data was released 11h 54m ago">00:00</td>
							<td class="left flagCur noWrap"><span title="Australia" class="ceFlags Australia" data-img_key="Australia">&nbsp;</span> AUD</td>
							<td class="left textNum sentiment noWrap" title="Moderate Volatility Expected" data-img_key="bull2"><i class="grayFullBullishIcon"></i><i class="grayFullBullishIcon"></i><i class="grayEmptyBullishIcon"></i></td>
							<td class="left event" title="Click to view more info on HIA New Home Sales"><a href="/economic-calendar/hia-new-home-sales-377" target="_blank">
							HIA New Home Sales (MoM) </a>
							</td>
							<td class="bold act blackFont event-436932-actual" title="" id="eventActual_436932">5.8%</td>
							<td class="fore  event-436932-forecast " id="eventForecast_436932">&nbsp;</td>
							<td class="prev blackFont  event-436932-previous" id="eventPrevious_436932"><span title="">-20.5%</span></td>
							<td class="alert js-injected-user-alert-container " data-name="HIA New Home Sales" data-event-id="377" data-status-enabled="0">
							<span class="js-plus-icon alertBellGrayPlus genToolTip oneliner" data-tooltip="Create Alert" data-tooltip-alt="Alert is active"></span>
							</td>
						</tr>
					</table>`,
			result: &InvestingScheduleRow{
				Id:           436932,
				EventId:      377,
				TimeStamp:    time.Date(2021, time.September, 16, 8, 0, 0, 0, time.UTC),
				CountryName:  "Australia",
				Title:        "HIA New Home Sales (MoM)",
				CurrencyCode: "AUD",
				Sentiment:    2,
				Actual:       &[]float64{5.8}[0],
				Forecast:     nil,
				Previous:     &[]float64{-20.5}[0],
				Type:         Index,
			},
			err: nil,
		},
		{
			html: `<table>
						<tr>
							<td colspan="9" class="theDay" id="theDay1631750400">Thursday, September 16, 2021</td>
						</tr>
					</table>`,
			result: nil,
			err:    nil,
		},
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		selector := goquery.NewDocumentFromNode(node).Find("tr")

		// Act
		parser := NewInvestingScheduleParser()
		value, err := parser.parseScheduleRowHtml(selector)

		// Assert
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.result, value)
	}
}

func TestParseScheduleRowId(t *testing.T) {
	tests := []struct {
		html   string
		result int
		err    error
	}{
		{
			html:   `<div id="row123456"></div>`,
			result: 123456,
			err:    nil,
		},
		{
			html:   `<div></div>`,
			result: 0,
			err:    &ParsingError{},
		},
		{
			html:   `<div id="rowABC"></div>`,
			result: 0,
			err:    &ParsingError{},
		},
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		selector := goquery.NewDocumentFromNode(node).Find("div")

		// Act
		parser := NewInvestingScheduleParser()
		value, err := parser.parseScheduleRowId(selector)

		// Assert
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.result, value)
	}
}

func TestParseScheduleCountryName(t *testing.T) {
	tests := []struct {
		html   string
		result string
		err    error
	}{
		{
			html:   `<tr><td><span class="ceFlags" title="China"></span></td></tr>`,
			result: "China",
			err:    nil,
		},
		{
			html:   `<tr><td><span class="ceFlags"></span></td></tr>`,
			result: "",
			err:    &ParsingError{},
		},
		{
			html:   `<tr><td></td><td></td><td></td></tr>`,
			result: "",
			err:    &ParsingError{},
		},
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		selector := goquery.NewDocumentFromNode(node).Selection

		// Act
		parser := NewInvestingScheduleParser()
		value, err := parser.parseScheduleCountryName(selector)

		// Assert
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.result, value)
	}
}

func TestParseScheduleTimeStamp(t *testing.T) {
	tests := []struct {
		html   string
		result time.Time
		err    error
	}{
		{
			html:   `<div data-event-datetime="2021/09/16 08:00:00"></div>`,
			result: time.Date(2021, time.September, 16, 8, 0, 0, 0, time.UTC),
			err:    nil,
		},
		{
			html:   `<div data-event-datetime="abcd"></div>`,
			result: time.Time{},
			err:    &ParsingError{},
		},
		{
			html:   `<div></div>`,
			result: time.Time{},
			err:    &ParsingError{},
		},
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		selector := goquery.NewDocumentFromNode(node).Find("div")

		// Act
		parser := NewInvestingScheduleParser()
		value, err := parser.parseScheduleTimeStamp(selector)

		// Assert
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.result, value)
	}
}

func TestParseScheduleEventType(t *testing.T) {
	tests := []struct {
		html   string
		result ScheduleEventType
		err    error
	}{
		{
			html: `<table>
						<tr>
							<td class="event"><a href="#">abc</a></td>
						</tr>
					</table>`,
			result: Index,
			err:    nil,
		},
		{
			html: `<table>
						<tr>
							<td class="event"><a href="#">abc</a> <span data-img_key="perliminary"></span></td>
						</tr>
					</table>`,
			result: PreliminaryRelease,
			err:    nil,
		},
		{
			html: `<table>
						<tr>
							<td class="event"><a href="#">abc</a> <span data-img_key="speach"></span></td>
						</tr>
					</table>`,
			result: Speech,
			err:    nil,
		},
		{
			html: `<table>
						<tr>
							<td class="event"><a href="#">abc</a> <span data-img_key="report"></span></td>
						</tr>
					</table>`,
			result: Report,
			err:    nil,
		},
		{
			html: `<table>
						<tr>
							<td class="event"><a href="#">abc</a> <span data-img_key="sandClock"></span></td>
						</tr>
					</table>`,
			result: RetrievingData,
			err:    nil,
		},
		{
			html: `<table>
						<tr>
							<td class="event"><a href="#">abc</a> <span></span></td>
						</tr>
					</table>`,
			result: Index,
			err:    &ParsingError{},
		},
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		selector := goquery.NewDocumentFromNode(node).Find("tr")

		// Act
		parser := NewInvestingScheduleParser()
		value, err := parser.parseScheduleEventType(selector)

		// Assert
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.result, value)
	}
}
