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

func convertHexdecimalInteger(content string) int64 {
	var value int64
	i := 0

	for i = 0; i < len(content); i++ {
		c := content[i]
		if c != '_' {
			value = (value << 4) | hexCharValues[c]
		}
	}

	return value
}

func ConvertHexdecimalInteger(content string) int64 {
	return convertHexdecimalInteger(content[2:])
}

func ConvertInteger(content string) int64 {
	if len(content) > 2 && content[:2] == "0x" {
		return ConvertHexdecimalInteger(content)
	}

	return ConvertDecimalInteger(content)
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

func convertDoubleQuoteString(content string) string {
	length := len(content) // In golang spec, len() returns the number of bytes.
	buffer := make([]byte, length)

	i := 1 // index of content, skip the first quote
	j := 0 // index of buffer
	finished := false

	for i < length && !finished {
		c := content[i]
		i++

		switch c {
		case '"':
			finished = true

		case '\\':
			e := content[i]
			i++
			switch e {
			case 'n':
				buffer[j] = '\n'
			case 'r':
				buffer[j] = '\r'
			case 't':
				buffer[j] = '\t'
			case '\\':
				buffer[j] = '\\'
			case '"':
				buffer[j] = '"'
			case 'x':
				code := content[i : i+2]
				buffer[j] = byte(convertHexdecimalInteger(code))
				i += 2
			}
			j++

		default:
			buffer[j] = c
			j++
		}
	}

	return string(buffer[:j])
}

func ConvertString(content string) string {
	var result string
	quote := content[0]

	switch quote {
	case '"':
		result = convertDoubleQuoteString(content)
	}

	return result
}
