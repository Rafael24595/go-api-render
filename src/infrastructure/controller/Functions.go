package controller

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Rafael24595/go-api-core/src/domain/body"
	"github.com/Rafael24595/go-api-core/src/domain/cookie"
	"github.com/google/uuid"
)

func Uuid() string {
	return uuid.NewString()
}

func Not(item any) bool {
	if item == nil {
		return true
	}
	if boolean, ok := item.(bool); ok {
		return !boolean
	}
	return false
}

func Concat(items ...string) string {
	return strings.Join(items, "")
}

func String(item any) string {
	if v, ok := item.([]byte); ok {
		return string(v)
	}
	return fmt.Sprintf("%v", item)
}

func BodyString(container string, payload body.Body) string {
	if container == constants.Body.TagText && payload.ContentType == body.Text {
		return String(payload.Bytes)
	}
	if container == constants.Body.TagJson && payload.ContentType == body.Json {
		return String(payload.Bytes)
	}
	return ""
}

func Join(items any, separator string) string {
	if sItems, ok := items.([]string); ok {
		return strings.Join(sItems, separator)
	}
	//TODO:
	return fmt.Sprintf("%v", items)
}

func FormatMilliseconds(ms int64) string {
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

func FormatMillisecondsDate(milliseconds int64) string {
	seconds := milliseconds / 1000
	nanoseconds := (milliseconds % 1000) * 1e6
	t := time.Unix(seconds, nanoseconds)
	return t.Format("2006-01-02 15:04:05")
}

func FormatBytes(bytes int) string {
	kb := float64(bytes) / 1024
	mb := kb / 1024
	gb := mb / 1024

	if round(gb) > 0 {
		return fmt.Sprintf("%.2f GB", gb)
	}
	if round(mb) > 0 {
		return fmt.Sprintf("%.2f MB", mb)
	}
	if round(kb) > 0 {
		return fmt.Sprintf("%.2f KB", kb)
	}

	return fmt.Sprintf("%.2f Bytes", float64(bytes))
}

func ParseCookie(cookie cookie.Cookie) string {
	return cookie.String()
}

func round(num float64) float64 {
	return math.Round(num)
}