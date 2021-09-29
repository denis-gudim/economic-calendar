package client

import (
	"economic-calendar/loader/investing/data"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type InvestingHtmlSourceMock struct {
	mock.Mock
}

func (mock *InvestingHtmlSourceMock) LoadEventsScheduleHtml(from, to time.Time, languageId int) (*goquery.Document, error) {

	if languageId == 2 {
		return nil, fmt.Errorf("test error")
	}

	html := `<table>
				<tr id="eventRowId_372">
					<td class="first left">All Day</td>
					<td class="flagCur left"><span title="South Korea" class="ceFlags South_Korea float_lang_base_1" data-img_key="South_Korea">&nbsp;</span></td>
					<td class="left textNum sentiment"><span class="bold">Holiday</span></td>
					<td colspan="6" class="left event">South Korea - Chuseok - Thanksgiving Day</td>
				</tr>
				<tr id="eventRowId_436019" class="js-event-item" event_attr_id="739" data-event-datetime="2021/09/20 02:00:00">
					<td class="first left time js-time" title="Event data will be released in 21h 59m">02:00</td>
					<td class="left flagCur noWrap"><span title="Germany" class="ceFlags Germany" data-img_key="Germany">&nbsp;</span> EUR</td>
					<td class="left textNum sentiment noWrap" title="Low Volatility Expected" data-img_key="bull1"><i class="grayFullBullishIcon"></i><i class="grayEmptyBullishIcon"></i><i class="grayEmptyBullishIcon"></i></td>
					<td class="left event" title="Click to view more info on German PPI"><a href="/economic-calendar/german-ppi-739" target="_blank">      German PPI (YoY)  (Aug)</a>      </td>
					<td class="bold act blackFont event-436019-actual" title="" id="eventActual_436019">&nbsp;</td>
					<td class="fore  event-436019-forecast " id="eventForecast_436019">11.4%</td>
					<td class="prev blackFont  event-436019-previous" id="eventPrevious_436019"><span title="">10.4%</span></td>
					<td class="alert js-injected-user-alert-container " data-name="German PPI" data-event-id="739" data-status-enabled="0">        <span class="js-plus-icon alertBellGrayPlus genToolTip oneliner" data-tooltip="Create Alert" data-tooltip-alt="Alert is active"></span>    </td>
				</tr>
				<tr id="eventRowId_437026" class="js-event-item" event_attr_id="559" data-event-datetime="2021/09/20 04:00:00">
					<td class="first left time js-time" title="Event data will be released in 23h 57m">04:00</td>
					<td class="left flagCur noWrap"><span title="Spain" class="ceFlags Spain" data-img_key="Spain">&nbsp;</span> EUR</td>
					<td class="left textNum sentiment noWrap" title="Low Volatility Expected" data-img_key="bull1"><i class="grayFullBullishIcon"></i><i class="grayEmptyBullishIcon"></i><i class="grayEmptyBullishIcon"></i></td>
					<td class="left event" title="Click to view more info on Spanish Trade Balance"><a href="/economic-calendar/spanish-trade-balance-559" target="_blank">      Spanish Trade Balance </a>      </td>
					<td class="bold act blackFont event-437026-actual" title="" id="eventActual_437026">&nbsp;</td>
					<td class="fore  event-437026-forecast " id="eventForecast_437026">&nbsp;</td>
					<td class="prev blackFont  event-437026-previous" id="eventPrevious_437026"><span title="">-0.98B</span></td>
					<td class="alert js-injected-user-alert-container " data-name="Spanish Trade Balance" data-event-id="559" data-status-enabled="0">        <span class="js-plus-icon alertBellGrayPlus genToolTip oneliner" data-tooltip="Create Alert" data-tooltip-alt="Alert is active"></span>    </td>
				</tr>
			</table>`
	return goquery.NewDocumentFromReader(strings.NewReader(html))
}

func (mocke *InvestingHtmlSourceMock) LoadEventDetailsHtml(eventId, languageId int) (*goquery.Document, error) {

	if languageId == 2 {
		return nil, fmt.Errorf("test error")
	}

	html := `
	<section id="leftColumn">
		<h1 class="ecTitle float_lang_base_1 relativeAttr">U.K. Core Retail Sales MoM	</h1>
		<div id="releaseInfo" class="releaseInfo bold">
			<span>Latest Release<div class="noBold">Sep 17, 2021</div></span>
			<span>Actual<div class="arial_14 redFont">71.0%</div></span>
			<span>Forecast<div class="arial_14 noBold">72.0%</div></span>
			<span>Previous<div class="arial_14 noBold blackFont">70.3%</div></span>
		</div>
		<div id="overViewBox" class="overViewBox event">
			<div class="left">The University of Michigan Consumer Sentiment Index.   </div>
			<div class="right">
				<div></div>
				<div></div>
				<div>
					<span>Source:</span>
					<span><a href="http://thomsonreuters.com/en/products-services/financial/investment-management.html" target="_blank" title="University of Michigan">University of Michigan</a></span>
				</div>
			</div>
		</div>
	</section>
	`
	return goquery.NewDocumentFromReader(strings.NewReader(html))
}

func (mock *InvestingHtmlSourceMock) LoadCountriesHtml(languageId int) (*goquery.Document, error) {

	if languageId == 2 {
		return nil, fmt.Errorf("test error")
	}

	html := `
		<div id="filtersWrapper">
			<ul class="countryOption">
				<li><input value="1"><label>Text 1</label></li>
				<li><input value="2"><label>Text 2</label></li>
			</ul>
		</div>`

	return goquery.NewDocumentFromReader(strings.NewReader(html))
}

func Test_InvestingRepository_getEventsScheduleByLanguage(t *testing.T) {
	// Arrange
	source := &InvestingHtmlSourceMock{}
	parser := &InvestingRepository{Source: source}
	time := time.Now()
	languageId := 1

	// Act
	actualResult, err := parser.getEventsScheduleByLanguage(languageId, time, time)

	// Assert
	source.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(actualResult))
	assert.Equal(t, languageId, actualResult[0].GetLanguageId())
	assert.Equal(t, languageId, actualResult[1].GetLanguageId())
}

func Test_InvestingRepository_GetEventsSchedule(t *testing.T) {
	// Arrange
	logger, hook := test.NewNullLogger()
	source := &InvestingHtmlSourceMock{}
	parser := &InvestingRepository{
		Source:            source,
		DefaultLanguageId: 1,
		BatchSize:         3,
		Logger:            logger,
	}
	time := time.Now()

	// Act
	actualResult, err := parser.GetEventsSchedule(time, time)

	// Assert
	source.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(actualResult))
	assert.Equal(t, (len(data.InvestingLanguagesMap) - 1), len(actualResult[436019]))
	assert.Equal(t, (len(data.InvestingLanguagesMap) - 1), len(actualResult[437026]))

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Contains(t, hook.LastEntry().Message, "test error")
}

func Test_InvestingRepository_getEventDetailsByLanguage(t *testing.T) {
	// Arrange
	source := &InvestingHtmlSourceMock{}
	parser := &InvestingRepository{Source: source}
	eventId := 123
	languageId := 1

	// Act
	actualResult, err := parser.getEventDetailsByLanguage(languageId, eventId)

	// Assert
	source.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(actualResult))
	assert.Equal(t, languageId, actualResult[0].GetLanguageId())
}

func Test_InvestingRepository_GetEventDetails(t *testing.T) {
	// Arrange
	logger, hook := test.NewNullLogger()
	source := &InvestingHtmlSourceMock{}
	parser := &InvestingRepository{
		Source:            source,
		DefaultLanguageId: 1,
		Logger:            logger,
	}
	eventId := 123

	// Act
	actualResult, err := parser.GetEventDetails(eventId)

	// Assert
	source.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, len(data.InvestingLanguagesMap)-1, len(actualResult))

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Contains(t, hook.LastEntry().Message, "test error")
}

func Test_InvestingRepository_getCountriesByLanguage(t *testing.T) {
	// Arrange
	source := &InvestingHtmlSourceMock{}
	parser := &InvestingRepository{Source: source}
	languageId := 1

	// Act
	actualResult, err := parser.getCountriesByLanguage(languageId)

	// Assert
	source.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(actualResult))
	assert.Equal(t, languageId, actualResult[0].GetLanguageId())
	assert.Equal(t, languageId, actualResult[1].GetLanguageId())
}

func Test_InvestingRepository_GetCountries(t *testing.T) {
	// Arrange
	logger, hook := test.NewNullLogger()
	source := &InvestingHtmlSourceMock{}
	parser := &InvestingRepository{
		Source:            source,
		DefaultLanguageId: 1,
		BatchSize:         4,
		Logger:            logger,
	}

	// Act
	actualResult, err := parser.GetCountries()

	// Assert
	source.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, len(data.InvestingLanguagesMap)-1, len(actualResult))
	assert.Equal(t, 2, len(actualResult[1]))

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Contains(t, hook.LastEntry().Message, "test error")
}
