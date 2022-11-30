package memory

import (
	"fmt"
	"strconv"
	"strings"
)

type UGeoTag string

func (g UGeoTag) GetLatitude() (float64, error) {
	arr := strings.Split(string(g), " ")
	if len(arr) < 2 {
		return 0, fmt.Errorf("invalid geotag %q", string(g))
	}

	val, err := strconv.ParseFloat(arr[0], 64)
	if err != nil {
		return 0, fmt.Errorf("parse geotag: %w", err)
	}

	return val, nil
}

func (g UGeoTag) GetLongitude() (float64, error) {
	arr := strings.Split(string(g), " ")
	if len(arr) < 2 {
		return 0, fmt.Errorf("invalid geotag %q", string(g))
	}

	val, err := strconv.ParseFloat(arr[1], 64)
	if err != nil {
		return 0, fmt.Errorf("parse geotag: %w", err)
	}

	return val, nil
}
