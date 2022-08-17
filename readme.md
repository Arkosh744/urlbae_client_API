# # API client for https://urlbae.com/developers

to install module:

    go get -u github.com/Arkosh744/urlbae_client_API

## Example
_________________________________________________
```go
package main

import (
	"fmt"
	"github.com/Arkosh744/urlbae_client_API/urlbae"
	"log"
	"time"
)

const apiToken = "PASTE_YOUR_API_TOKEN_HERE"

func main() {
	// Create a new client
	client, err := urlbae.NewClient(apiKey, time.Second*5)
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
	}

	// Get account info and check login status
	message, err := urlbae.GetAccountInfo(client)
	log.Print("Account info:\n" +
		"Email: " + message.Email + "\n" +
		"Username: " + message.Username + "\n" +
		"Status: " + message.Status + "\n" +
		"Registered at: " + message.RegisteredAt + "\n")

	// Write down your link data
	longUrlname := "https://www.yandex.ru"
	customName := "yandex"
	expirationDate := time.Now().Add(time.Hour * 24) // set to 1 day from now

	UrlToShort := &urlbae.LongLinkData{LongURL: longUrlname, CustomName: customName, ExpirationDate: expirationDate}

	// Shorten your url
	urlbae.DoShortLink(client, UrlToShort)

	// Get List of all active short links
	_, err = urlbae.GetAllLinks(client)
	if err != nil {
		log.Println(err)
	}
}
```
_________________________________________________

