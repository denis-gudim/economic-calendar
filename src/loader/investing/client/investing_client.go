package client

import (
	"compress/gzip"
	"economic-calendar/loader/investing/data"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
)

type InvestingHttpClient struct {
	RetryCount int
}

func (client *InvestingHttpClient) LoadEventDetailsHtml(eventId, languageId int) (*goquery.Document, error) {
	language := data.InvestingLanguagesMap[languageId]
	url := fmt.Sprintf("https://%s.investing.com/economic-calendar/%x-%d", language.Domain, [16]byte(uuid.New()), eventId)

	return client.doHtmlRequest("GET", url, nil, nil)
}

func (client *InvestingHttpClient) LoadEventsScheduleHtml(from, to time.Time, languageId int) (response *goquery.Document, err error) {

	language := data.InvestingLanguagesMap[languageId]
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

	responseJson, err := client.doJsonRequest("POST", requestUrl, &headers, &params)

	if err != nil {
		return
	}

	data, ok := responseJson["data"]

	if !ok {
		err = fmt.Errorf("invalid response JSON. data property not found")
		return
	}

	tbl := fmt.Sprintf("<table>%s</table>", data)

	return goquery.NewDocumentFromReader(strings.NewReader(tbl))
}

func (client *InvestingHttpClient) LoadCountriesHtml(languageId int) (*goquery.Document, error) {
	language := data.InvestingLanguagesMap[languageId]
	url := fmt.Sprintf("https://%s.investing.com/economic-calendar/?_uid=%x", language.Domain, [16]byte(uuid.New()))

	return client.doHtmlRequest("GET", url, nil, nil)
}

func (client *InvestingHttpClient) doJsonRequest(method, url string, headers *http.Header, body *url.Values) (response map[string]interface{}, err error) {
	reader, err := client.doRetryRequest(method, url, headers, body)

	if err != nil {
		return
	}

	defer reader.Close()

	err = json.NewDecoder(reader).Decode(&response)

	return
}

func (client *InvestingHttpClient) doHtmlRequest(method, url string, headers *http.Header, body *url.Values) (response *goquery.Document, err error) {
	reader, err := client.doRetryRequest(method, url, headers, body)

	if err != nil {
		return
	}

	defer reader.Close()

	return goquery.NewDocumentFromReader(reader)
}

func (client *InvestingHttpClient) doRetryRequest(method, url string, headers *http.Header, body *url.Values) (reader *gzip.Reader, err error) {

	for i := 0; i < client.RetryCount; i++ {
		reader, err = client.doRequest(method, url, headers, body)

		if err == nil {
			break
		}

		if reader != nil {
			reader.Close()
		}
	}

	return
}

func (client *InvestingHttpClient) doRequest(method, url string, headers *http.Header, body *url.Values) (reader *gzip.Reader, err error) {

	var bodyReader io.Reader

	if body != nil {
		bodyReader = strings.NewReader(shuffleRequestParams(body))
	}

	request, err := http.NewRequest(method, url, bodyReader)

	if err != nil {
		return
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
		return
	}

	responseEncoding := response.Header.Get("Content-Encoding")

	if responseEncoding != "gzip" {
		err = fmt.Errorf("invalid response encoding '%s'", responseEncoding)
		return
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
