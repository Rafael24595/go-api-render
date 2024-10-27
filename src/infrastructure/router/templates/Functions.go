package templates

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func not(item any) bool {
	if item == nil {
		return true
	}
	if boolean, ok := item.(bool); ok {
		return !boolean
	}
	if str, ok := item.(string); ok {
		return str == ""
	}
	return false
}

func itemString(item any) string {
	if v, ok := item.([]byte); ok {
		return string(v)
	}
	return fmt.Sprintf("%v", item)
}

func concat(items ...string) string {
	return strings.Join(items, "")
}

func join(items any, separator string) string {
	if sItems, ok := items.([]string); ok {
		return strings.Join(sItems, separator)
	}
	return fmt.Sprintf("%v", items)
}

func uuidString() string {
	return uuid.NewString()
}

func millisecondsToTime(ms int64) string {
	duration := time.Duration(ms) * time.Millisecond

	hours := int64(duration.Hours())
	minutes := int64(duration.Minutes()) % 60
	seconds := int64(duration.Seconds()) % 60
	milliseconds := ms % 1000

	if hours > 0 {
		return fmt.Sprintf("%dh %dm :%ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	if seconds > 0 {
		return fmt.Sprintf("%ds", seconds)
	}

	return fmt.Sprintf("%dms", milliseconds)
}

func millisecondsToDate(milliseconds int64) string {
	seconds := milliseconds / 1000
	nanoseconds := (milliseconds % 1000) * 1e6
	t := time.Unix(seconds, nanoseconds)
	return t.Format("2006-01-02 15:04:05")
}