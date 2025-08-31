package wbgt

import (
	"fmt"
)

// AlertLevel represents the heat stroke alert level as an enum-like type
type AlertLevel int8

// Alert level constants - these are the only valid values for AlertLevel
const (
	NoAlert         AlertLevel = 0
	HeatStrokeAlert AlertLevel = 1
	SevereAlert     AlertLevel = 3
	OutOfTime       AlertLevel = 9
)

// String implements the Stringer interface for better printing
func (a AlertLevel) String() string {
	switch a {
	case NoAlert:
		return "NoAlert"
	case HeatStrokeAlert:
		return "HeatStrokeAlert"
	case SevereAlert:
		return "HeatStrokeAlert - Severe"
	case OutOfTime:
		return "OutOfTime"
	default:
		return fmt.Sprintf("Unknown(%d)", int(a))
	}
}

// IsAlert returns true if there's an active heat stroke alert
func (a AlertLevel) IsAlertActive() bool {
	return a == HeatStrokeAlert || a == SevereAlert
}

type CSVRecord struct {
	Prefecture     string     // Column 1: 府県予報区
	PrefectureCode string     // Column 6: 都道府県コード
	AlertFlag      AlertLevel // Column 7: TargetDate1フラグ
	WBGT           string     // Column 11: 日最高WBGT（5:00）
}

// Constants
const (
	TokyoLocationTitle = "東京都"
	TokyoCityName      = "東京"
	TokyoLocationCode  = 13
	AlertThreshold     = 33.0
)

// IsHeatStrokeAlert checks if the data indicates a heat stroke alert
func (w *CSVRecord) IsHeatStrokeAlert() bool {
	return w.AlertFlag.IsAlertActive()
}

// GetAlertLevel returns the alert level based on WBGT value
func (w *CSVRecord) GetAlertLevel() AlertLevel {
	if w.AlertFlag.IsAlertActive() {
		return w.AlertFlag
	}

	return NoAlert
}
