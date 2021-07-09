package poloniex

import (
	"encoding/json"
	"strconv"
	"time"
)

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func parseJSONFloatString(data json.RawMessage) (float64, error) {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(s, 64)
}

func parseStringToTime(t string) (time.Time, error) {
	// "2021-07-09 03:46:50"
	year, err := strconv.Atoi(t[:4])
	if err != nil {
		return time.Time{}, Error(WrongTimeFormat)
	}
	month, err := strconv.Atoi(t[5:7])
	if err != nil {
		return time.Time{}, Error(WrongTimeFormat)
	}
	day, err := strconv.Atoi(t[8:10])
	if err != nil {
		return time.Time{}, Error(WrongTimeFormat)
	}
	hours, err := strconv.Atoi(t[11:13])
	if err != nil {
		return time.Time{}, Error(WrongTimeFormat)
	}
	minutes, err := strconv.Atoi(t[14:16])
	if err != nil {
		return time.Time{}, Error(WrongTimeFormat)
	}
	seconds, err := strconv.Atoi(t[17:])
	if err != nil {
		return time.Time{}, Error(WrongTimeFormat)
	}

	return time.Date(year, time.Month(month), day, hours, minutes, seconds, 0, time.UTC), nil
}
