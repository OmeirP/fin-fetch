package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
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

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)	// maybe replace with logging
		os.Exit(1)
	}
	defer resp.Body.Close()	 // release socket when surrounding func (main) is finished. If body not fully read, connection is closed and not returned to pool.
}