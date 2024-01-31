package utils

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

func FetchData(apiURL string) ([]byte, error) {
	res, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch data due to (%v)", err)
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body (%v)", err)
	}

	return data, nil
}

func IsValidUUID(id string) bool {
	uuidPattern := "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89aAbB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$"
	regex := regexp.MustCompile(uuidPattern)

	return regex.MatchString(id)
}

func OnlyLetters(name string) bool {
	regex := regexp.MustCompile("^[a-zA-Z]+$")
	return regex.MatchString(name)
}
