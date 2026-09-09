package main

import (
	"encoding/base64"
	"encoding/json"
	"encoding/csv\"
	"github.com/joho/godotenv"
	"fmt"
	"io"
	"log"
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
	PPL float64	// Price Profit/Loss	its a slash not a divide

}


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	apiKeyId := os.Getenv("T212_KEY_ID")
    apiKeySecret := os.Getenv("T212_KEY_SECRET")

	url := "https://live.trading212.com/api/v0/equity/portfolio"


	credentials := apiKeyId + ":" + apiKeySecret
	encodedCred := base64.StdEncoding.EncodeToString([]byte(credentials))

	authHeader := "Basic " + encodedCred

	// Prepare request in memory. http.Get instead would instantly send req without letting change header
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)	// 1 is normally handled error
	}


	req.Header.Set("Authorization", authHeader)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer resp.Body.Close()	 // release socket when surrounding func (main) is finished. If body not fully read, connection is closed and not returned to pool.


	// resp.Body is a readCloser so can use ReadAll
	bodyBytes, err := io.ReadAll(resp.Body)	// bodyBytes is json-encoded bytes
	if err != nil {
		log.Fatal(err)
	}
	

	// an array of the above defined position structs
	var positions []Position

	// parses response into the position slices. Each slice is an instance of the struct. The values get placed into the appropriate fields automatically.
	err = json.Unmarshal(bodyBytes, &positions)
	if err != nil {
		log.Fatal(err)
	}


	// iterate through the range of positions. Use _ if you don't need the index
	for _, pos := range positions {
		fmt.Println(pos)
	}
}