package main

import (
	"context"
	"fmt"

	"heat-alert-bot/internal/wbgt"
)

func main() {
	// Generate endpoint for current day
	endpoint := wbgt.GetAlertEndpoint()

	client := wbgt.NewClient(wbgt.HeatAlertDownloadURL)
	ctx := context.Background()

	data, err := client.FetchCSVData(ctx, endpoint)
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return
	}

	fmt.Printf("Fetched %d bytes of CSV data\n", len(data))
}
