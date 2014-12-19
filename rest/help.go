package rest

import (
	"strings"
)

//pattern  url pattern,such as a/b/c ; a/b/*/d
//url      a/b/c/
func match(pattern, url string) bool {
	reqParts := strings.Split(url, "/")
	patParts := strings.Split(pattern, "/")
	i, j := 0, 0
	for i < len(reqParts) && j < len(patParts) {
		if reqParts[i] == patParts[j] {
			i++
			j++
		} else if patParts[j] == "*" {
			if j+1 < len(patParts) && patParts[j+1] == reqParts[i] {
				j++
				j++
			}
			i++

		} else if strings.Contains(patParts[j], "{") {
			i++
			j++
		} else {
			break
		}
	}
	if i == len(reqParts) && j == len(patParts) {
		return true
	}

	return false

}
