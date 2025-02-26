package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Config struct {
    SpreadsheetID   string
    SheetName       string
    CredentialsFile string
}

func main() {
    // Configuration
    config := Config{
        SpreadsheetID:   "YOUR_SPREADSHEET_ID",
        SheetName:       "Sheet1",
        CredentialsFile: "credentials.json",
    }

    // Initialize context and client
    ctx := context.Background()
    srv, err := sheets.NewService(ctx, option.WithCredentialsFile(config.CredentialsFile))
    if err != nil {
        log.Fatalf("Unable to retrieve Sheets client: %v", err)
    }

    // Get initial row count
    rangeName := fmt.Sprintf("%s!A:C", config.SheetName) // A:C covers timestamp, email, name
    resp, err := srv.Spreadsheets.Values.Get(config.SpreadsheetID, rangeName).Do()
    if err != nil {
        log.Fatalf("Unable to retrieve initial data from sheet: %v", err)
    }

    // Set initial row count and log it
    lastRowCount := len(resp.Values)
    log.Printf("Program started. Initial row count: %d", lastRowCount)

    // Poll the spreadsheet every 5 seconds
    for {
        // Get spreadsheet data
        resp, err := srv.Spreadsheets.Values.Get(config.SpreadsheetID, rangeName).Do()
        if err != nil {
            log.Printf("Unable to retrieve data from sheet: %v", err)
            time.Sleep(5 * time.Second)
            continue
        }

        // Check if there are new rows
        if len(resp.Values) > lastRowCount {
            newRows := resp.Values[lastRowCount:]
            for _, row := range newRows {
                if len(row) >= 3 { // Ensure we have timestamp, email, and name
                    timestamp := row[0].(string)
                    email := row[1].(string)
                    name := row[2].(string)
                    
                    // Print to console for new rows only
                    fmt.Printf("Email would be sent - Timestamp: %s, Email: %s, Name: %s\n",
                        timestamp, email, name)
                    log.Printf("Processed new response for %s (%s)", name, email)
                }
            }
            lastRowCount = len(resp.Values)
        }

        time.Sleep(5 * time.Second) // Wait before next poll
    }
}