package handlers

import (
	"net/http"
	"strconv"
)

func GetQueryParamString(r *http.Request, key string) (string, bool) {
	keys := r.URL.Query()
	result := keys.Get(key)
	if result == "" {
		return "", false
	}
	return result, true
}

func GetQueryParamInt(r *http.Request, key string) (int, bool) {
	keys := r.URL.Query()
	result := keys.Get(key)
	i, err := strconv.Atoi(result)
	if err != nil {
		return 0, false
	}
	return i, true
}

func GetDefaultQueryParams(r *http.Request) (int, int) {
	limit, ok := GetQueryParamInt(r, "limit")
	if !ok {
		limit = 100
	}
	page, ok := GetQueryParamInt(r, "page")
	if !ok || page < 1 {
		page = 1
	}
	return limit, page
}
