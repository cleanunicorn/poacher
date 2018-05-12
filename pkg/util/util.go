package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// GetJSON downloads the JSON and runs unmarshal on the provided interface
func GetJSON(url string, object interface{}) error {
	responseBody, err := GetURLBody(url)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, &object)
	if err != nil {
		return err
	}

	return nil
}

// GetURLBody downloads the URL and returns it
func GetURLBody(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return []byte{}, err
	}

	return responseBody, nil
}
