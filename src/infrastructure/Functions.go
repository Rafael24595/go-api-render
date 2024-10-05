package infrastructure

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func ToString(value any) string {
	return fmt.Sprintf("%v", value)
}

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