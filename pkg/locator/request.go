package locator

import (
	"io/ioutil"
	"net/http"
)

func GET(URL string) ([]byte, error) {
	r, err := http.Get(URL)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
