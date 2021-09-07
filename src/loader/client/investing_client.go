package client

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

type InvestingHttpClient struct {
	RetryCount int
}

func (client *InvestingHttpClient) LoadEventDetailsHtml(eventId int, languageId int) (string, error) {
	language := InvestingLanguagesMap[int32(languageId)]
	url := fmt.Sprintf("https://%s.investing.com/economic-calendar/%x-%d", language.domain, [16]byte(uuid.New()), eventId)

	response, err := client.doRetryRequest("GET", url, nil, nil)

	if err == nil {
		return string(response), nil
	}

	return "", err
}

func (client *InvestingHttpClient) LoadEventsScheduleHtml(from time.Time, to time.Time, languageId int) (string, error) {

	language := InvestingLanguagesMap[int32(languageId)]
	refererUrl := fmt.Sprintf("https://%s.investing.com/economic-calendar", language.domain)
	requestUrl := fmt.Sprintf("%s/Service/getCalendarFilteredData", refererUrl)

	headers := http.Header{
		"Accept":  {"application/json, text/javascript, */*; q=0.01"},
		"Referer": {refererUrl},
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
		"dateFrom":      {from.Format("yyyy-MM-dd")},
		"dateTo":        {to.Format("yyyy-MM-dd")},
		"uuid":          {uuid.New().String()},
	}

	response, err := client.doRetryRequest("POST", requestUrl, headers, params)

	if err != nil {
		return "", err
	}

	var bodyJson map[string]interface{}

	if err := json.Unmarshal(response, &bodyJson); err != nil {
		return "", err
	}

	return bodyJson["data"].(string), nil
}

func (client *InvestingHttpClient) LoadCountriesHtml(languageId int) (string, error) {
	language := InvestingLanguagesMap[int32(languageId)]
	url := fmt.Sprintf("https://%s.investing.com/economic-calendar/?_uid=%x", language.domain, [16]byte(uuid.New()))

	response, err := client.doRetryRequest("GET", url, nil, nil)

	if err == nil {
		return string(response), nil
	}

	return "", err
}

func (client *InvestingHttpClient) doRetryRequest(method string, url string, headers http.Header, body url.Values) ([]byte, error) {
	var err error

	for i := 0; i < client.RetryCount; i++ {
		resp, err := doRequest(method, url, headers, body)

		if err == nil {
			return resp, err
		} else {
			fmt.Println(err)
			continue
		}
	}

	return nil, err
}

func doRequest(method string, url string, headers http.Header, body url.Values) ([]byte, error) {

	var bodyReader io.Reader

	if body != nil {
		bodyReader = strings.NewReader(shuffleRequestParams(body))
	}

	request, err := http.NewRequest(method, url, bodyReader)

	if err != nil {
		return nil, err
	}

	if headers != nil {
		request.Header = headers
	} else {
		request.Header.Set("Accept", "*/*")
	}

	request.Header.Set("Accept-Encoding", "gzip")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36")
	request.Header.Set("X-Requested-With", "XMLHttpRequest")

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, err
	}

	var reader io.ReadCloser

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		if reader, err = gzip.NewReader(response.Body); err != nil {
			return nil, err
		}
	default:
		reader = response.Body
	}

	defer reader.Close()

	return ioutil.ReadAll(reader)
}

func shuffleRequestParams(body url.Values) string {
	params := strings.Split(body.Encode(), "&")
	perm := rand.Perm(len(params))

	for i, v := range perm {
		params[i], params[v] = params[v], params[i]
	}

	return strings.Join(params, "&")
}
