package helpers

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"
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

func PaginationParse(r *http.Request) (int, int, error) {
	const max_limit = 200

	var limit, offset int = 100, 0
	var err error

	limitString := r.FormValue("limit")
	offsetString := r.FormValue("offset")

	if limitString != "" {
		limit, err = strconv.Atoi(limitString)
		if err != nil {
			return 0, 0, err
		}
	}
	if offsetString != "" {
		offset, err = strconv.Atoi(offsetString)
		if err != nil {
			return 0, 0, err
		}
	}

	if limit <= 0 || offset < 0 || limit > max_limit {
		return 0, 0, errors.New("0<limit<=" + strconv.FormatInt(max_limit, 10) + "&& offset > 0")
	}

	return limit, offset, nil
}

func GetQueryArgs(args [][]string, r *http.Request) map[string]interface{} {
	return getQueryArgs(args, r, false)
}

func GetQueryArgsStrict(args [][]string, r *http.Request) map[string]interface{} {
	return getQueryArgs(args, r, true)
}

func getQueryArgs(args [][]string, r *http.Request, argsMustExist bool) map[string]interface{} {
	filter := make(map[string]interface{})
	for _, names := range args {
		namesLength := len(names)
		switch namesLength {
		case 1:
			inputName := names[0]
			value := r.FormValue(inputName)
			if value != "" {
				filter[inputName] = value
			} else if argsMustExist {
				return nil
			}
		case 2:
			inputName := names[0]
			dbName := names[1]
			value := r.FormValue(inputName)
			if value != "" {
				filter[dbName] = value
			} else if argsMustExist {
				return nil
			}
		default:
			return nil
		}
	}
	return filter
}

func GetContentType(rawContentType string) string {
	if rawContentType == "" {
		return rawContentType
	}
	return strings.Split(rawContentType, ";")[0]
}

func GetContentLength(rawContentLength string) int64 {
	if rawContentLength == "" {
		return -1
	}
	contentLength, err := strconv.ParseInt(rawContentLength, 10, 64)
	if err != nil {
		return -1
	}
	return contentLength
}

func ParseInt64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

func ParseInt(value string) (int, error) {
	integer, err := strconv.ParseInt(value, 10, 0)
	return int(integer), err
}

func Int64ToString(n int64) string {
	return strconv.FormatInt(n, 10)
}

func FastConcat(strings ...string) string {
	var buffer bytes.Buffer
	for i := range strings {
		buffer.WriteString(strings[i])
	}
	return buffer.String()
}
