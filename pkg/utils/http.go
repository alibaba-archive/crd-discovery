package utils

import "net/http"

func GetStringOrElse(r *http.Request, key string, el string) string {
	values, ok := r.URL.Query()[key]
	if !ok || len(values) < 1 {
		return el
	}
	return values[0]
}
