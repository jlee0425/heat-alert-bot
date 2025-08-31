package main

import (
	"fmt"
	"os"

	"heat-alert-bot/internal/wbgt"
)

func main() {
	// Test with sample data first
	fmt.Println("------- Testing with sample.csv -------")
	sampleData, err := os.ReadFile("sample.csv")
	if err != nil {
		fmt.Printf("Error reading sample.csv: %v\n", err)
	} else {
		testParser(sampleData, "sample.csv")
	}
}

func testParser(data []byte, source string) {
	fmt.Printf("Testing parser with data from %s:\n", source)

	IsAlertActive, wbgtTemp, alertLevel, err := wbgt.CheckTokyoHeatAlert(data)
	if err != nil {
		fmt.Printf("Error checking Tokyo heat alert: %v\n", err)
		return
	}

	fmt.Printf("Tokyo Heat Alert Status:\n")
	fmt.Printf("   - Alert Active: %t\n", IsAlertActive)
	fmt.Printf("   - WBGT Temperature: %.1f°C\n", wbgtTemp)
	fmt.Printf("   - Alert Level: %s\n", alertLevel)

	if IsAlertActive {
		fmt.Printf("HEAT STROKE ALERT: Remote work recommended\n")
	} else {
		fmt.Printf("No heat alert: Normal work conditions\n")
	}
}
