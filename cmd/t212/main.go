package main

import (
	"encoding/base64"
	"encoding/json"
	"encoding/csv"
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
	AveragePrice float64
	CurrentPrice float64
	PPL float64	// Price Profit/Loss	its a slash not a divide

}


const (
	csvName = "portfolio_data.csv"
	dateFormat = "02/01/06"
)


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


	// CSV

	// if the file doesn't exist, create it, otherwise append.  order of flags doesn't matter.
	// Append puts write cursor at end, create just checks if it exists and creates if not - does nothing if so.
	f, err := os.OpenFile("trade_log.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, )	// O_WRONLY is required by POSIX, open requests must declare intent.
	if err != nil {
		log.Fatal(err)	// This is if there's an error opening the file.
	}
	defer f.Close()


	if _, err := f.Write([]byte("appended some data\n")); err != nil {
		log.Fatal(err)
	}


}



func updateCSV(currentPositions []Position) {
	now := time.Now()
	currDateStr := now.Format(dateFormat)
	cutoff := now.AddDate(-10, 0, 0)		// Making cutoff date 10 years ago.

	headers := []string{"SnapshotDate", "Ticker", "Quantity", "AveragePrice", "CurrentPrice", "PPL"}



	// Create preserved rows and then write the headings at the top.
	preservedRows := [][]string{headers}	// 2d slice but its basically a 2d array	
	preservedRows = append(preservedRows, headers)
	
	
	// Read existing csv and get rid of old records if necessary
	// done with this scope because file doesn't need to be read or exist for the program. If it fails, it should just move on and create a new one.
	// also doesn't need to be used after this is done.
	if file, err := os.Open(csvName); err == nil {
		reader := csv.NewReader(file)

		// read/move past header
		_. _ = reader.Read()

		for {	// empty for is basically a while true
			record, err := reader.Read()

			if err == io.EOF {
				break
			}

			if err != nil || len(record) < 1 {
				continue
			}

			rowDate, err := time.Parse(dateFormat, record[0])
			if err != nil {
				continue  // skip if something wrong with rows.
			}

			// raw record as string checked instead of parsing date to avoid needing to truncate time from now value
			if (rowDate.After(cutoff) || rowDate.Equal(cutoff)) && record[0] != currDateStr {
				preservedRows = append(preservedRows, record)
			}
		}
		file.Close()
	}


	// Now append currently fetched snapshot
	for _, p := range currentPositions {
		currValue := p.CurrentPrice * p.Quantity
		row := []string{
			todayStr,
			p.Ticker,
			fmt.Sprintf("%.8f", p.Quantity),
			fmt.Sprintf("%.2f", p.AveragePrice),
			fmt.Sprintf("%.2f", p.CurrentPrice),
			fmt.Sprintf("%.2f", currValue),
			fmt.Sprintf("%.2f", p.PPL),
		}
		preservedRows = append(preservedRows, row)
	}


	// Create new temp file and update that one first. Replace the actual file when writing is complete in case of crash.
	tmpFileName := csvName + ".tmp"
	tmpFile, err := os.Create(tmpFileName)
	if err != nil {
		log.Fatal("Couldn't create temp file:", err)
	}

	writer := csv.NewWriter(tmpFile)
	if err := writer.WriteAll(preservedRows); err != nil {
		tmpFile.Close()
		os.Remove(tmpFileName)	// Remove the file on error
	}
}
