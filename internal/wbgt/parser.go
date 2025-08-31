package wbgt

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

func parseWBGTCSV(data []byte) ([]CSVRecord, error) {
	reader := csv.NewReader(strings.NewReader(string(data)))

	var records []CSVRecord
	var headerPassed bool

	for {
		row, err := reader.Read()
		if err != nil {
			break // End of file or error
		}

		// Skip header rows until we find the data section
		// Data section starts with "府県予報区" header
		if !headerPassed {
			if len(row) > 0 && strings.Contains(row[0], "府県予報区") {
				headerPassed = true
			}
			continue
		}

		// Skip empty rows
		if len(row) == 0 || strings.TrimSpace(row[0]) == "" {
			continue
		}

		// Parse data row (need at least 11 columns)
		if len(row) < 11 {
			continue
		}

		// Parse alert flag (column 7, index 6)
		alertFlag := NoAlert
		if len(row) > 6 {
			if flagValue, err := strconv.Atoi(strings.TrimSpace(row[6])); err == nil {
				alertFlag = AlertLevel(flagValue)
			}
		}

		record := CSVRecord{
			Prefecture:     strings.TrimSpace(row[0]),  // Column 1
			PrefectureCode: strings.TrimSpace(row[5]),  // Column 6 (index 5)
			AlertFlag:      alertFlag,                  // Column 7 (index 6)
			WBGT:           strings.TrimSpace(row[10]), // Column 11 (index 10)
		}

		records = append(records, record)
	}

	return records, nil
}

// FindTokyoAlert finds Tokyo's alert information from the records
func findTokyoAlert(records []CSVRecord) (*CSVRecord, error) {
	for _, record := range records {
		// Tokyo can be identified by prefecture name or code
		if record.Prefecture == TokyoLocationTitle || record.PrefectureCode == "13" {
			return &record, nil
		}
	}
	return nil, fmt.Errorf("tokyo data not found in CSV")
}

// ParseTokyoWBGT extracts Tokyo's WBGT value from the WBGT data string
// Format: "小河内:31/青梅:33/練馬:33/八王子:33/府中:33/東京:33/江戸川臨海:33/..."
func parseTokyoWBGT(wbgtData string) (float64, error) {
	// Split by "/" to get individual location:temperature pairs
	locations := strings.SplitSeq(wbgtData, "/")

	for location := range locations {
		// Split by ":" to get location and temperature
		parts := strings.Split(location, ":")
		if len(parts) != 2 {
			continue
		}

		locationName := strings.TrimSpace(parts[0])
		tempStr := strings.TrimSpace(parts[1])

		// Look for Tokyo (東京)
		if locationName == TokyoCityName {
			temp, err := strconv.ParseFloat(tempStr, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse Tokyo temperature: %w", err)
			}
			return temp, nil
		}
	}

	return 0, fmt.Errorf("tokyo temperature not found in WBGT data")
}

// CheckTokyoHeatAlert checks if there's a heat alert for Tokyo
func CheckTokyoHeatAlert(data []byte) (bool, float64, AlertLevel, error) {
	records, err := parseWBGTCSV(data)

	for row := range records {
		fmt.Println(records[row])
	}

	if err != nil {
		return false, 0, NoAlert, fmt.Errorf("failed to parse CSV: %w", err)
	}

	tokyoRecord, err := findTokyoAlert(records)
	if err != nil {
		return false, 0, NoAlert, err
	}

	tokyoWBGT, err := parseTokyoWBGT(tokyoRecord.WBGT)
	if err != nil {
		return false, 0, NoAlert, err
	}

	// Check if alert is active (flag = 1 or 3) or temperature >= threshold
	IsAlertActive := tokyoRecord.AlertFlag.IsAlertActive() || tokyoWBGT >= AlertThreshold

	return IsAlertActive, tokyoWBGT, tokyoRecord.AlertFlag, nil
}
