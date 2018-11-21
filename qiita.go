package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type QiitaAPI struct {
	Domain      string
	AccessToken string
	PerPage     int
}

type QiitaItem struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	User      struct {
		ID              string `json:"id"`
		ProfileImageURL string `json:"profile_image_url"`
	} `json:"user"`
}

func (item *QiitaItem) Time_CreatedAt() (time.Time, error) {
	// "created_at": "2000-01-01T00:00:00+00:00",
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Qiita: <error>: %v\n", err)
		return time.Now(), err
	}

	t, err := time.Parse(time.RFC3339, item.CreatedAt)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Qiita: <error>: %v\n", err)
		return time.Now(), err
	}

	return t.In(jst), nil
}

func (item *QiitaItem) dateDescription() string {
	t, err := item.Time_CreatedAt()
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func (item *QiitaItem) String() string {
	return fmt.Sprintf("%v %v %v", item.Title, item.CreatedAt, item.URL)
}

func (qiita *QiitaAPI) requestURLString(minDate string, maxDate string) string {
	baseURL := fmt.Sprintf("https://%s/api/v2/items", qiita.Domain)
	query := fmt.Sprintf("created:>%s created:<%s", minDate, maxDate)

	parameters := url.Values{}
	parameters.Add("query", query)
	parameters.Add("per_page", "50")

	queryString := strings.Replace(parameters.Encode(), "+", "%20", -1)
	queryString = strings.Replace(queryString, ":", "%3A", -1)

	return baseURL + "?" + queryString
}

func (qiita *QiitaAPI) Items(minDate string, maxDate string) ([]*QiitaItem, error) {
	url := qiita.requestURLString(minDate, maxDate)

	// Prepare HTTP Request
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", "Bearer "+qiita.AccessToken)

	// Fetch Request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Qiita: <error> fetch stocks:", err)
		return nil, err
	}

	// Read Response Body
	responseBody, _ := ioutil.ReadAll(response.Body)

	if response.Status != "200 OK" {
		fmt.Println("response Status : ", response.Status)
		fmt.Println("response Headers : ", response.Header)
		fmt.Println("response Body : ", string(responseBody))
		return nil, nil
	}

	// Decode JSON
	var items []*QiitaItem
	if err := json.Unmarshal(responseBody, &items); err != nil {
		fmt.Fprintln(os.Stderr, "Qiita: <error> json unmarshal:", err)
		return nil, err
	}

	return items, nil
}

func NewQiitaAPI(config Config) *QiitaAPI {
	qiita := QiitaAPI{
		Domain:      config.Qiita.Domain,
		AccessToken: config.Qiita.AccessToken,
		PerPage:     config.Qiita.PerPage,
	}

	if qiita.PerPage == 0 {
		qiita.PerPage = 5
	}

	return &qiita
}
