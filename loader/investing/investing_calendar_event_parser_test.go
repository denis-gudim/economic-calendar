package investing

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func Test_InvestingCalendarEventParser_ParseCalendarEventHtml(t *testing.T) {
	tests := []struct {
		html           string
		expectedResult *InvestingCalendarEvent
		err            error
	}{
		{
			html: `<section id="leftColumn">
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
								<div>
									<span>
										<i class="grayFullBullishIcon"></i>
										<i class="grayFullBullishIcon"></i>
										<i class="grayEmptyBullishIcon"></i>
									</span>
								</div>
								<div>
									<span>
										<i title="Japan" class="ceFlags Japan middle inlineblock"></i>
									</span>
								</div>
								<div>
									<span>Source:</span>
									<span><a href="http://thomsonreuters.com/en/products-services/financial/investment-management.html" target="_blank" title="University of Michigan">University of Michigan</a></span>
								</div>
							</div>
						</div>
					</section>`,
			expectedResult: &InvestingCalendarEvent{
				Title:     "U.K. Core Retail Sales MoM",
				Overview:  "The University of Michigan Consumer Sentiment Index.",
				Source:    "University of Michigan",
				SourceUrl: "http://thomsonreuters.com/en/products-services/financial/investment-management.html",
				Unit:      "%",
				Sentiment: 2,
				Country:   "Japan",
			},
			err: nil,
		},
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		html := goquery.NewDocumentFromNode(node)
		parser := NewInvestingCalendarEventParser()

		// Act
		actualResult, err := parser.ParseCalendarEventHtml(html)

		// Assert
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expectedResult, actualResult)
	}
}

func Test_InvestingCalendarEventParser_parserTitle(t *testing.T) {
	tests := []struct {
		html           string
		expectedResult string
		err            error
	}{
		{
			html: `<section id="leftColumn">
						<h1 class="ecTitle float_lang_base_1 relativeAttr">U.K. Core Retail Sales MoM	</h1>
					</section>`,
			expectedResult: "U.K. Core Retail Sales MoM",
			err:            nil,
		},
		{
			html: `<section id="leftColumn">
						<h1 class="ecTitle float_lang_base_1 relativeAttr"></h1>
					</section>`,
			expectedResult: "",
			err:            fmt.Errorf("parsing error"),
		},
		{
			html: `<section id="leftColumn">
						<h1>U.K. Core Retail Sales MoM	</h1>
					</section>`,
			expectedResult: "",
			err:            fmt.Errorf("parsing error"),
		},
		{
			html:           `<section id="leftColumn"></section>`,
			expectedResult: "",
			err:            fmt.Errorf("parsing error"),
		},
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		selection := goquery.NewDocumentFromNode(node).Find("#leftColumn")
		parser := NewInvestingCalendarEventParser()

		// Act
		actualResult, err := parser.parseTitle(selection)

		// Assert
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expectedResult, actualResult)
	}
}

func Test_InvestingCalendarEventParser_ParseUnit(t *testing.T) {
	tests := []struct {
		html           string
		expectedResult string
		err            error
	}{
		{
			html: `<section id="leftColumn">
						<div id="releaseInfo" class="releaseInfo bold">
							<span>Latest Release<div class="noBold">Sep 17, 2021</div></span>
							<span>Actual<div class="arial_14 redFont">-71.0B</div></span>
							<span>Forecast<div class="arial_14 noBold">&nbsp;</div></span>
							<span>Previous<div class="arial_14 noBold blackFont">70.3B</div></span>
						</div>
					</section>`,
			expectedResult: "B",
			err:            nil,
		},
		{
			html: `<section id="leftColumn">
						<div id="releaseInfo" class="releaseInfo bold">
							<span>Latest Release<div class="noBold">Sep 17, 2021</div></span>
							<span>Actual<div class="arial_14 redFont">7.1%</div></span>
							<span>Forecast<div class="arial_14 noBold">&nbsp;</div></span>
							<span>Previous<div class="arial_14 noBold blackFont">7.1%</div></span>
						</div>
					</section>`,
			expectedResult: "%",
			err:            nil,
		},
		{
			html: `<section id="leftColumn">
						<div id="releaseInfo" class="releaseInfo bold">
							<span>Latest Release<div class="noBold">Sep 17, 2021</div></span>
							<span>Actual<div class="arial_14 redFont">-37.2</div></span>
							<span>Forecast<div class="arial_14 noBold">&nbsp;</div></span>
							<span>Previous<div class="arial_14 noBold blackFont">-35.8</div></span>
						</div>
					</section>`,
			expectedResult: "",
			err:            nil,
		},
		{
			html: `<section id="leftColumn">
						<div id="releaseInfo" class="releaseInfo bold">
						</div>
					</section>`,
			expectedResult: "",
			err:            nil,
		},
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		selection := goquery.NewDocumentFromNode(node).Find("#releaseInfo")
		parser := NewInvestingCalendarEventParser()

		// Act
		actualResult, err := parser.parseUnit(selection)

		// Assert
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expectedResult, actualResult)
	}
}

func Test_InvestingCalendarEventParser_parseSentiment(t *testing.T) {
	tests := []struct {
		html           string
		expectedResult int
		err            error
	}{
		{
			html: `<div id="overViewBox">
						<div class="right">
							<span>
								<i class="grayFullBullishIcon"></i>
								<i class="grayEmptyBullishIcon"></i>
								<i class="grayEmptyBullishIcon"></i>
							</span>
						</div>
					</div>`,
			expectedResult: 1,
			err:            nil,
		},
		{
			html: `<div id="overViewBox">
						<div class="right">
							<span></span>
						</div>
					</div>`,
			expectedResult: 0,
			err:            fmt.Errorf("parse error"),
		},
		{
			html: `<div id="overViewBox">
						<div>
						</div>
					</div>`,
			expectedResult: 0,
			err:            fmt.Errorf("parse error"),
		},
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		selection := goquery.NewDocumentFromNode(node).Find("#overViewBox")
		parser := NewInvestingCalendarEventParser()

		// Act
		actualResult, err := parser.parseSentiment(selection)

		// Assert
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expectedResult, actualResult)
	}
}

func Test_InvestingCalendarEventParser_parseCountry(t *testing.T) {
	tests := []struct {
		html           string
		expectedResult string
		err            error
	}{
		{
			html: `<div id="overViewBox">
						<div class="right">
							<span>
								<i title="Japan" class="ceFlags Japan middle inlineblock"></i>
							</span>
						</div>
					</div>`,
			expectedResult: "Japan",
			err:            nil,
		},
		{
			html: `<div id="overViewBox">
						<div class="right">
							<span></span>
						</div>
					</div>`,
			expectedResult: "",
			err:            fmt.Errorf("parse error"),
		},
	}

	for _, test := range tests {
		// Arrange
		node, _ := html.Parse(strings.NewReader(test.html))
		selection := goquery.NewDocumentFromNode(node).Find("#overViewBox")
		parser := NewInvestingCalendarEventParser()

		// Act
		actualResult, err := parser.parseCountry(selection)

		// Assert
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expectedResult, actualResult)
	}
}
