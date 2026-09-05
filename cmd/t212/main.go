package main

import (
	"encoding/json"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
)

// typed data template
type Position struct {
	// capital field names so encoding/json pkg can read and populate the field
	// field names here don't need the matching tags since case insensitive matching is used by default
	Ticker string
	Quantity float64
	CurrentPrice float64
	PPL float64

}


func main() {
	apiKeyId := os.Getenv("T212_KEY_ID")
    apiKeySecret := os.Getenv("T212_KEY_SECRET")

	url := "https://live.trading212.com/api/v0/equity/portfolio"


	credentials := apiKeyId + ":" + apiKeySecret
	encoded_cred := base64.StdEncoding.EncodeToString([]byte(credentials))

	auth_header := "Basic " + encoded_cred

	// Prepare request in memory. http.Get instead would instantly send req without letting change header
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)	// 1 is normally handled error
	}


	req.Header.Set("Authorization", "Basic " + encoded_cred)
}