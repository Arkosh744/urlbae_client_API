package urlbae

import (
	"bytes"
	"encoding/json"
	"errors"
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
	log.Println("Shortening the url: " + LongLinkData.LongURL + " with the alias: " + LongLinkData.CustomName + " and expiration date: " + LongLinkData.ExpirationDate.Format("2006-01-02 15:04:05"))
	ShortedLink, _ := shortLink(client, LongLinkData)

	// If alias exists - we generate new random alias
	if ShortedLink.Message == CustomNameExists {
		log.Println("Sorry, that alias is taken. We generated random one.")
		LongLinkData.CustomName = ""
		log.Println("Shortening the url: " + LongLinkData.LongURL + " without alias and expiration date: " + LongLinkData.ExpirationDate.Format("2006-01-02 15:04:05"))
		ShortedLink, _ = shortLink(client, LongLinkData)
	}
	log.Println("New shorted link is: " + ShortedLink.ShortUrl)
	return ShortedLink.ShortUrl
}

func GetAccountInfo(client *client) (AccountInfo, error) {
	url := "https://urlbae.com/api/account"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return AccountInfo{}, err
	}

	req.Header = http.Header{
		"Authorization": {"Bearer " + client.apiKey},
		"Content-Type":  {"application/json"},
	}
	res, err := client.client.Do(req)
	if err != nil {
		return AccountInfo{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing body: %s", err)
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return AccountInfo{}, err
	}

	AccountDataResponse := LinkResponse{}
	err = json.Unmarshal(body, &AccountDataResponse)
	if err != nil {
		return AccountInfo{}, err
	}

	if AccountDataResponse.Message == InvalidKeyAPI {
		return AccountInfo{}, errors.New("Invalid API key")
	}

	return AccountDataResponse.AccountInfo, nil
}

func shortLink(client *client, LongLinkData *LongLinkData) (LinkResponse, error) {
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
	if err != nil {
		return LinkResponse{}, err
	}

	req.Header = http.Header{
		"Authorization": {"Bearer " + client.apiKey},
		"Content-Type":  {"application/json"},
	}
	res, err := client.client.Do(req)
	if err != nil {
		return LinkResponse{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing body: %s", err)
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return LinkResponse{}, err
	}

	ShortLinkResponseData := LinkResponse{}

	err = json.Unmarshal(body, &ShortLinkResponseData)
	if err != nil {
		return LinkResponse{}, err
	}

	return ShortLinkResponseData, nil
}

func GetAllLinks(client *client) ([]GeneratedLinkData, error) {
	url := "https://urlbae.com/api/urls?order=date"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Authorization": {"Bearer " + client.apiKey},
		"Content-Type":  {"application/json"},
	}
	res, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing body: %s", err)
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	AllLinksData := ListLinksResponse{}
	err = json.Unmarshal(body, &AllLinksData)
	for _, link := range AllLinksData.LinkData.URLS {
		log.Printf(link.ShortURL + " for " + link.LongURL + " Total clicks:" + link.Clicks)
	}
	return AllLinksData.LinkData.URLS, nil

}
