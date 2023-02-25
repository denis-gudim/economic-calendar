package investing

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func Test_InvestingScheduleParser_ParseScheduleRowHtml(t *testing.T) {
	tests := []struct {
		html           string
		expectedResult *InvestingScheduleRow
		err            error
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
			expectedResult: &InvestingScheduleRow{
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
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		selector := goquery.NewDocumentFromNode(node).Find("tr")

		// Act
		parser := NewInvestingScheduleParser()
		actualResult, err := parser.parseScheduleRowHtml(selector)

		// Assert
		assert.Equal(t, test.err, err)
		assert.Equal(t, test.expectedResult, actualResult)
	}
}

func Test_InvestingScheduleParser_ParseScheduleRowId(t *testing.T) {
	tests := []struct {
		html           string
		expectedResult int
		err            error
	}{
		{
			html:           `<div id="row123456"></div>`,
			expectedResult: 123456,
			err:            nil,
		},
		{
			html:           `<div></div>`,
			expectedResult: 0,
			err:            fmt.Errorf("parse error"),
		},
		{
			html:           `<div id="rowABC"></div>`,
			expectedResult: 0,
			err:            fmt.Errorf("parse error"),
		},
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		selector := goquery.NewDocumentFromNode(node).Find("div")

		// Act
		parser := NewInvestingScheduleParser()
		actualResult, err := parser.parseScheduleRowId(selector)

		// Assert
		assert.Equal(t, test.err, err)
		assert.Equal(t, test.expectedResult, actualResult)
	}
}

func Test_InvestingScheduleParser_ParseScheduleCountryName(t *testing.T) {
	tests := []struct {
		html           string
		expectedResult string
		err            error
	}{
		{
			html:           `<tr><td><span class="ceFlags" title="China"></span></td></tr>`,
			expectedResult: "China",
			err:            nil,
		},
		{
			html:           `<tr><td><span class="ceFlags"></span></td></tr>`,
			expectedResult: "",
			err:            fmt.Errorf("parse error"),
		},
		{
			html:           `<tr><td></td><td></td><td></td></tr>`,
			expectedResult: "",
			err:            fmt.Errorf("parse error"),
		},
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		selector := goquery.NewDocumentFromNode(node).Selection

		// Act
		parser := NewInvestingScheduleParser()
		actualResult, err := parser.parseScheduleCountryName(selector)

		// Assert
		assert.Equal(t, test.err, err)
		assert.Equal(t, test.expectedResult, actualResult)
	}
}

func Test_InvestingScheduleParser_ParseScheduleTimeStamp(t *testing.T) {
	tests := []struct {
		html           string
		expectedResult time.Time
		err            error
	}{
		{
			html:           `<div data-event-datetime="2021/09/16 08:00:00"></div>`,
			expectedResult: time.Date(2021, time.September, 16, 8, 0, 0, 0, time.UTC),
			err:            nil,
		},
		{
			html:           `<div data-event-datetime="abcd"></div>`,
			expectedResult: time.Time{},
			err:            fmt.Errorf("parse error"),
		},
		{
			html:           `<div></div>`,
			expectedResult: time.Time{},
			err:            fmt.Errorf("parse error"),
		},
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		selector := goquery.NewDocumentFromNode(node).Find("div")

		// Act
		parser := NewInvestingScheduleParser()
		actualResult, err := parser.parseScheduleTimeStamp(selector)

		// Assert
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expectedResult, actualResult)
	}
}

func Test_InvestingScheduleParser_ParseScheduleEventType(t *testing.T) {
	tests := []struct {
		html           string
		expectedResult ScheduleEventType
		err            error
	}{
		{
			html: `<table>
						<tr>
							<td class="event"><a href="#">abc</a></td>
						</tr>
					</table>`,
			expectedResult: Index,
			err:            nil,
		},
		{
			html: `<table>
						<tr>
							<td class="event"><a href="#">abc</a> <span data-img_key="perliminary"></span></td>
						</tr>
					</table>`,
			expectedResult: PreliminaryRelease,
			err:            nil,
		},
		{
			html: `<table>
						<tr>
							<td class="event"><a href="#">abc</a> <span data-img_key="speach"></span></td>
						</tr>
					</table>`,
			expectedResult: Speech,
			err:            nil,
		},
		{
			html: `<table>
						<tr>
							<td class="event"><a href="#">abc</a> <span data-img_key="report"></span></td>
						</tr>
					</table>`,
			expectedResult: Report,
			err:            nil,
		},
		{
			html: `<table>
						<tr>
							<td class="event"><a href="#">abc</a> <span data-img_key="sandClock"></span></td>
						</tr>
					</table>`,
			expectedResult: RetrievingData,
			err:            nil,
		},
		{
			html: `<table>
						<tr>
							<td class="event"><a href="#">abc</a> <span></span></td>
						</tr>
					</table>`,
			expectedResult: Index,
			err:            fmt.Errorf("parse error"),
		},
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		selector := goquery.NewDocumentFromNode(node).Find("tr")

		// Act
		parser := NewInvestingScheduleParser()
		actualResult, err := parser.parseScheduleEventType(selector)

		// Assert
		assert.Equal(t, test.err, err)
		assert.Equal(t, test.expectedResult, actualResult)
	}
}
