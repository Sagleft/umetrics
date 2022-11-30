package locator

import (
	"encoding/json"
	"fmt"
)

type defaultLocator struct {
}

func NewDefaultLocator() Locator {
	return &defaultLocator{}
}

func (l *defaultLocator) GetCoordinates(address string) (Coordinates, error) {
	response, err := GET(fmt.Sprintf(
		"http://ipwho.is/%s",
		address,
	))
	if err != nil {
		return Coordinates{}, err
	}

	c := Coordinates{}
	err = json.Unmarshal(response, &c)
	return c, err
}
