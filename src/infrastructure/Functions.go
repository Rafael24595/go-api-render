package infrastructure

import (
	"fmt"

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