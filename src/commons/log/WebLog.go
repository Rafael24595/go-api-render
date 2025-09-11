package log

import (
	"bytes"

	"github.com/Rafael24595/go-api-core/src/commons/log"
)

type WebLog struct {
}

func NewWebLog() WebLog {
	return WebLog{}
}

func (w WebLog) Custome(category string, err error) {
	log.Custome(category, err)
}

func (w WebLog) Custom(category string, message string) {
	log.Custom(category, message)
}

func (w WebLog) Customf(category string, format string, args ...any) {
	log.Customf(category, format, args...)
}

func (w WebLog) Message(message string) {
	log.Message(message)
}

func (w WebLog) Messagef(format string, args ...any) {
	log.Messagef(format, args...)
}

func (w WebLog) Warning(message string) {
	log.Warning(message)
}

func (w WebLog) Warningf(format string, args ...any) {
	log.Warningf(format, args...)
}

func (w WebLog) Error(err error) {
	log.Error(err)
}

func (w WebLog) Errors(message string) {
	log.Errors(message)
}

func (w WebLog) Errorf(format string, args ...any) {
	log.Errorf(format, args...)
}

func (w WebLog) Write(slice []byte) (int, error) {
	w.Warningf("%s", bytes.TrimSpace(slice))
	return len(slice), nil
}