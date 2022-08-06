package urlbae

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type client struct {
	// The HTTP client
	client *http.Client
	apiKey string
}

func NewClient(apiKey string, timeout time.Duration) (*client, error) {

	if timeout == 0 {
		return nil, errors.New("timeout can't be 0")
	}

	return &client{
		client: &http.Client{Timeout: timeout,
			Transport: &loggingRoundTripper{logger: os.Stdout, next: http.DefaultTransport}},
		apiKey: apiKey,
	}, nil
}

func DoShortLink(client *client, LongLinkData *LongLinkData) string {
	fmt.Println("Shortening the url: " + LongLinkData.LongURL + " with the alias: " + LongLinkData.CustomName + " and expiration date: " + LongLinkData.ExpirationDate.Format("2006-01-02 15:04:05"))
	ShortedLink := shortLink(client, LongLinkData)

	if ShortedLink.Message == CustomNameExists {
		fmt.Println("Sorry, that alias is taken. We generated random one.")
		LongLinkData.CustomName = ""
		fmt.Println("Shortening the url: " + LongLinkData.LongURL + " without alias and expiration date: " + LongLinkData.ExpirationDate.Format("2006-01-02 15:04:05"))
		ShortedLink = shortLink(client, LongLinkData)
	}
	fmt.Println("Your new shorted link is: " + ShortedLink.ShortUrl)
	return ShortedLink.ShortUrl
}

func GetAccountInfo(client *client) LinkResponse {
	url := "https://urlbae.com/api/account"

	req, err := http.NewRequest("GET", url, nil)
	ErrorFatalEnding(err)

	req.Header = http.Header{
		"Authorization": {"Bearer " + client.apiKey},
		"Content-Type":  {"application/json"},
	}

	res, err := client.client.Do(req)
	ErrorFatalEnding(err)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		ErrorFatalEnding(err)
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	ErrorFatalEnding(err)

	AccountDataResponse := LinkResponse{}
	err = json.Unmarshal(body, &AccountDataResponse)
	ErrorFatalEnding(err)

	return AccountDataResponse
}

func shortLink(client *client, LongLinkData *LongLinkData) LinkResponse {
	url := "https://urlbae.com/api/url/add"

	// if expirationDate is not set, it will be set to 1 day from now
	if LongLinkData.ExpirationDate == (time.Time{}) {
		LongLinkData.ExpirationDate = time.Now().Add(time.Hour * 24)
	}

	data := []byte(`{"url": "` + LongLinkData.LongURL + `", "custom": "` + LongLinkData.CustomName + `", "expiry": "` + LongLinkData.ExpirationDate.Format("2006-01-02 15:04:05") + `"}`)

	if LongLinkData.CustomName == "" {
		data = []byte(`{"url": "` + LongLinkData.LongURL + `", "expiry": "` + LongLinkData.ExpirationDate.Format("2006-01-02 15:04:05") + `"}`)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	ErrorFatalEnding(err)

	req.Header = http.Header{
		"Authorization": {"Bearer " + client.apiKey},
		"Content-Type":  {"application/json"},
	}

	res, err := client.client.Do(req)
	ErrorFatalEnding(err)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		ErrorFatalEnding(err)
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	ErrorFatalEnding(err)

	ShortLinkResponseData := LinkResponse{}

	err = json.Unmarshal(body, &ShortLinkResponseData)
	ErrorFatalEnding(err)

	return ShortLinkResponseData
}

func GetAllLinks(client *client) {
	url := "https://urlbae.com/api/urls?order=date"

	req, err := http.NewRequest("GET", url, nil)
	ErrorFatalEnding(err)

	req.Header = http.Header{
		"Authorization": {"Bearer " + client.apiKey},
		"Content-Type":  {"application/json"},
	}

	res, err := client.client.Do(req)
	ErrorFatalEnding(err)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		ErrorFatalEnding(err)
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	ErrorFatalEnding(err)

	AllLinksData := ListLinksResponse{}

	err = json.Unmarshal(body, &AllLinksData)

	for _, link := range AllLinksData.LinkData.URLS {
		log.Printf(link.ShortURL + " for " + link.LongURL + " Total clicks:" + link.Clicks)
	}

}

func ErrorFatalEnding(err error) {
	if err != nil {
		//Handle Error
		log.Fatalln(err)
	}
}
