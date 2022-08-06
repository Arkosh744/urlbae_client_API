package main

import (
	"shortURL_client_API/urlbae"
	"time"
)

const apiToken = "PASTE_YOUR_API_TOKEN_HERE"

func main() {
	// Create a new client
	client, err := urlbae.NewClient(apiToken, time.Second*5)
	urlbae.ErrorFatalEnding(err)

	// Get account info
	urlbae.GetAccountInfo(client)

	// Create a new long link data
	longUrlname := "https://www.google.ru"
	customName := "goodsagle"
	expirationDate := time.Now().Add(time.Hour * 24) // set to 1 day from now

	UrlToShort := &urlbae.LongLinkData{longUrlname, customName, expirationDate}

	// Shorten the url
	urlbae.DoShortLink(client, UrlToShort)

	// Get List of all active short links
	urlbae.GetAllLinks(client)

}
