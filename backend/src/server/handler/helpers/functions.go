package helpers

import (
	"bytes"
	"strings"
)

func Error(error string) string {
	escapedError, err := escapeJsonString(error)
	if err != nil {
		escapedError = "there was an error escaping the error string for json"
	}
	return `{"error":"` + escapedError + `"}`
}

// Adapted from google code html.EscapeString
func escapeJsonString(s string) (string, error) {
	const escapedChars = `"`

	i := strings.IndexAny(s, escapedChars)
	if i == -1 {
		return s, nil
	}

	var buf bytes.Buffer
	for i != -1 {
		if _, err := buf.WriteString(s[:i]); err != nil {
			return "", err
		}
		var esc string
		switch s[i] {
		case '"':
			esc = `\"`
		default:
			panic("unrecognized escape character")
		}
		s = s[i+1:]
		if _, err := buf.WriteString(esc); err != nil {
			return "", err
		}
		i = strings.IndexAny(s, escapedChars)
	}
	_, err := buf.WriteString(s)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
