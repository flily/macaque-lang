package parser

import (
	"strconv"
)

func ConvertDecimalInteger(content string) int64 {
	var value int64

	for i := range content {
		c := content[i]
		if c >= '0' && c <= '9' {
			value = (value * 10) + int64(c-'0')
		}
	}

	return value
}

var hexCharValues = [...]int64{
	'0': 0,
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'A': 10, 'a': 10,
	'B': 11, 'b': 11,
	'C': 12, 'c': 12,
	'D': 13, 'd': 13,
	'E': 14, 'e': 14,
	'F': 15, 'f': 15,
}

func ConvertHexdecimalInteger(content string) int64 {
	var value int64
	i := 0

	for i = 2; i < len(content); i++ {
		c := content[i]
		if c != '_' {
			value = (value << 4) | hexCharValues[c]
		}
	}

	return value
}

func ConvertFloat(content string) float64 {
	buffer := make([]byte, 0, len(content))
	for i := range content {
		c := content[i]
		if c != '_' {
			buffer = append(buffer, c)
		}
	}

	number := string(buffer)
	result, _ := strconv.ParseFloat(number, 64)
	return result
}

func ConvertString(content string) string {
	buffer := make([]byte, 0, len(content))

	return string(buffer)
}
