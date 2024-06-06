package auth

import (
	"errors"
	"net/http"
	"strings"
)

var ErrNoAuthHeaderIncluded = errors.New("no authorizatoin header included")

func GetSentApiKey(header http.Header) (string, error) {
	sentApiKey := header.Get("Authorization")
	if sentApiKey == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitKey := strings.Split(sentApiKey, " ")

	if len(splitKey) < 2 || splitKey[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}
	return splitKey[1], nil

}
