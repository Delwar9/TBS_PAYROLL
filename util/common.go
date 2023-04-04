package util

import (
	"fmt"
	"time"
)

// GetFinancialYear Returns the current financial year
func GetFinancialYear() string {
	y, m, d := time.Now().Date()
	if m >= time.July && d >= 1 {
		return fmt.Sprintf("%d-%d", y, y+1)
	}
	return fmt.Sprintf("%d-%d", y-1, y)
}

// GetFinancialYearFromDate returns financial year from input date
func GetFinancialYearFromDate(date string) (string, error) {
	dt, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", fmt.Errorf("invalid date provided")
	}
	y, m, d := dt.Date()
	if m >= time.July && d >= 1 {
		return fmt.Sprintf("%d-%d", y, y+1), nil
	}
	return fmt.Sprintf("%d-%d", y-1, y), nil
}

// GetFinancialYearFromDateTime returns financial year from input date
func GetFinancialYearFromDateTime(date *time.Time) string {
	y, m, d := date.Date()
	if m >= time.July && d >= 1 {
		return fmt.Sprintf("%d-%d", y, y+1)
	}
	return fmt.Sprintf("%d-%d", y-1, y)
}
