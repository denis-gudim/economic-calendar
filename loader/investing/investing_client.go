package investing

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/denis-gudim/economic-calendar/loader"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
)

type InvestingHttpClient struct {
	RetryCount int
}

func NewInvestingHttpClient(cnf *loader.Config) *InvestingHttpClient {
	return &InvestingHttpClient{
		RetryCount: cnf.Loading.RetryCount,
	}
}

func (client *InvestingHttpClient) LoadEventDetailsHtml(ctx context.Context, eventId, languageId int) (*goquery.Document, error) {
	language := InvestingLanguagesMap[languageId]
	url := fmt.Sprintf("https://%s.investing.com/economic-calendar/%x-%d", language.Domain, [16]byte(uuid.New()), eventId)

	return client.doHtmlRequest(ctx, "GET", url, nil, nil)
}

func (client *InvestingHttpClient) LoadEventsScheduleHtml(ctx context.Context, from, to time.Time, languageId int) (response *goquery.Document, err error) {

	language := InvestingLanguagesMap[languageId]
	refererUrl := fmt.Sprintf("https://%s.investing.com/economic-calendar", language.Domain)
	requestUrl := fmt.Sprintf("%s/Service/getCalendarFilteredData", refererUrl)

	headers := http.Header{
		"Accept":           {"application/json, text/javascript, */*; q=0.01"},
		"Content-Type":     {"application/x-www-form-urlencoded"},
		"Referer":          {refererUrl},
		"X-Requested-With": {"XMLHttpRequest"},
	}

	params := url.Values{
		"country[]":     {"29", "25", "54", "145", "47", "34", "174", "163", "32", "70", "6", "232", "27", "37", "122", "15", "78", "113", "107", "55", "24", "121", "59", "89", "72", "71", "22", "17", "51", "39", "93", "106", "14", "48", "66", "33", "23", "10", "119", "35", "92", "102", "57", "94", "97", "68", "96", "103", "111", "42", "109", "188", "7", "139", "247", "105", "172", "21", "43", "20", "60", "87", "44", "193", "125", "45", "53", "38", "170", "100", "56", "80", "52", "238", "36", "90", "112", "110", "11", "26", "162", "9", "12", "46", "85", "41", "202", "63", "123", "61", "143", "4", "5", "138", "178", "84", "75"},
		"category[]":    {"_employment", "_economicActivity", "_inflation", "_credit", "_centralBanks", "_confidenceIndex", "_balance", "_Bonds"},
		"importance[]":  {"1", "2", "3"},
		"timeZone":      {"55"},
		"timeFilter":    {"timeOnly"},
		"currentTab":    {"custom"},
		"submitFilters": {"1"},
		"limit_from":    {"0"},
		"dateFrom":      {from.Format("2006-01-02")},
		"dateTo":        {to.Format("2006-01-02")},
		"uuid":          {uuid.New().String()},
	}

	responseJson, err := client.doJsonRequest(ctx, "POST", requestUrl, &headers, &params)

	if err != nil {
		return
	}

	data, ok := responseJson["data"]

	if !ok {
		err = fmt.Errorf("investing client do request: invalid response JSON data property not found")
		return
	}

	tbl := fmt.Sprintf("<table>%s</table>", data)

	return goquery.NewDocumentFromReader(strings.NewReader(tbl))
}

func (client *InvestingHttpClient) LoadCountriesHtml(ctx context.Context, languageId int) (*goquery.Document, error) {
	language := InvestingLanguagesMap[languageId]
	url := fmt.Sprintf("https://%s.investing.com/economic-calendar/?_uid=%x", language.Domain, [16]byte(uuid.New()))

	return client.doHtmlRequest(ctx, "GET", url, nil, nil)
}

func (client *InvestingHttpClient) doJsonRequest(ctx context.Context, method, url string, headers *http.Header, body *url.Values) (response map[string]interface{}, err error) {
	reader, err := client.doRetryRequest(ctx, method, url, headers, body)

	if err != nil {
		return
	}

	defer reader.Close()

	err = json.NewDecoder(reader).Decode(&response)

	return
}

func (client *InvestingHttpClient) doHtmlRequest(ctx context.Context, method, url string, headers *http.Header, body *url.Values) (response *goquery.Document, err error) {
	reader, err := client.doRetryRequest(ctx, method, url, headers, body)

	if err != nil {
		return
	}

	defer reader.Close()

	return goquery.NewDocumentFromReader(reader)
}

func (client *InvestingHttpClient) doRetryRequest(ctx context.Context, method, url string, headers *http.Header, body *url.Values) (reader *gzip.Reader, err error) {

	for i := 0; i < client.RetryCount; i++ {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("do retry request canceled")
		default:
			{
				reader, err = client.doRequest(ctx, method, url, headers, body)

				if err == nil {
					break
				}

				if reader != nil {
					reader.Close()
				}
			}
		}
	}

	return
}

func (client *InvestingHttpClient) doRequest(ctx context.Context, method, url string, headers *http.Header, body *url.Values) (reader *gzip.Reader, err error) {

	var bodyReader io.Reader

	if body != nil {
		bodyReader = strings.NewReader(shuffleRequestParams(body))
	}

	request, err := http.NewRequestWithContext(ctx, method, url, bodyReader)

	if err != nil {
		return nil, fmt.Errorf("investing client create request: %w", err)
	}

	if headers != nil {
		request.Header = *headers
	} else {
		request.Header.Set("Accept", "*/*")
	}

	request.Header.Set("Accept-Encoding", "gzip")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36")

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, fmt.Errorf("investing client do request error: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("investing client do request: invalid response code '%v'", response.StatusCode)
	}

	encoding := response.Header.Get("Content-Encoding")

	if encoding != "gzip" {
		return nil, fmt.Errorf("investing client do request: invalid response encoding '%v'", encoding)
	}

	return gzip.NewReader(response.Body)
}

func shuffleRequestParams(body *url.Values) string {
	params := strings.Split(body.Encode(), "&")
	perm := rand.Perm(len(params))

	for i, v := range perm {
		params[i], params[v] = params[v], params[i]
	}

	return strings.Join(params, "&")
}
