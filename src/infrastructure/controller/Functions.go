package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"

	"github.com/Rafael24595/go-api-core/src/domain/body"
	"github.com/Rafael24595/go-api-core/src/domain/cookie"
)

func BodyString(container string, payload body.Body) string {
	if container == constants.Body.TagText && payload.ContentType == body.Text {
		return string(payload.Bytes)
	}
	if container == constants.Body.TagJson && payload.ContentType == body.Json {
		return string(payload.Bytes)
	}
	return ""
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

func FormatXml(input string) string {
	//TODO: implement
	return input
}

func FormatHtml(input string) string {
	//TODO: implement
	return input
}

func FormatJson(input string) string {
	var prettyJson bytes.Buffer
	err := json.Indent(&prettyJson, []byte(input), "", "    ")
	if err != nil {
		println(err.Error())
		prettyJson = *bytes.NewBuffer([]byte(input))
	}
	return prettyJson.String()
}

func round(num float64) float64 {
	return math.Round(num)
}
