package urlbae

import "time"

type LinkResponse struct {
	Error       int         `json:"error,omitempty"`
	Id          string      `json:"id,omitempty"`
	ShortUrl    string      `json:"shorturl,omitempty"`
	Message     string      `json:"message,omitempty"`
	AccountInfo AccountInfo `json:"data,omitempty"`
}

type AccountInfo struct {
	AccountId    int    `json:"id,omitempty"`
	Email        string `json:"email,omitempty"`
	Username     string `json:"username,omitempty"`
	Status       string `json:"status,omitempty"`
	RegisteredAt string `json:"registered,omitempty"`
}

type LongLinkData struct {
	LongURL        string
	CustomName     string
	ExpirationDate time.Time
}

type ListLinksResponse struct {
	LinkData struct {
		URLS []GeneratedLinkData `json:"urls"`
	} `json:"data"`
}

type GeneratedLinkData struct {
	Id             string `json:"id"`
	Alias          string `json:"alias"`
	ShortURL       string `json:"shorturl"`
	LongURL        string `json:"longurl"`
	Clicks         string `json:"clicks"`
	Title          string `json:"title"`
	ExpirationDate string `json:"date"`
}

const CustomNameExists = "That alias is taken. Please choose another one."
const InvalidKeyAPI = "A valid API key is required to use this service."
