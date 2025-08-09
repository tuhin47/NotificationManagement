package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateRandomNumber(i int) string {
	if i <= 0 {
		return ""
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	min := 1
	for k := 0; k < i-1; k++ {
		min *= 10
	}
	max := min*10 - 1

	if i == 1 {
		min = 0
	}

	return strconv.Itoa(r.Intn(max-min+1) + min)
}

func Contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func RecurrenceDuration(every uint, recurrence string, current *time.Time) time.Duration {
	var minutes uint
	switch recurrence {
	case "seconds":
		return time.Duration(every) * time.Second
	case "minutes":
		return time.Duration(every) * time.Minute
	case "hour":
		return time.Duration(every) * time.Hour
	case "daily":
		return time.Duration(every) * time.Hour * 24
	case "weekly":
		return time.Duration(every) * time.Hour * 24 * 7
	case "monthly":
		return AddMonthDate(current, 1).Sub(*current)
	case "quarterly":
		return AddMonthDate(current, 4).Sub(*current)
	default:
		minutes = 0
	}
	return time.Duration(minutes) * time.Second
}

func AddMonthDate(t *time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}
