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
	client, err := urlbae.NewClient(apiToken, time.Second*5)
	urlbae.ErrorFatalEnding(err)

	// Get account info and check login status
	message := urlbae.GetAccountInfo(client)
	if message.Message == urlbae.InvalidKeyAPI {
		log.Fatalln("Invalid API key")
	} else {
		fmt.Print("Account info:\n" +
			"Email: " + message.AccountInfo.Email + "\n" +
			"Username: " + message.AccountInfo.Username + "\n" +
			"Status: " + message.AccountInfo.Status + "\n" +
			"Registered at: " + message.AccountInfo.RegisteredAt + "\n")
	}

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
```
_________________________________________________

