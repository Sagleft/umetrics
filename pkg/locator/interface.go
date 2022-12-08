package locator

type Locator interface {
	GetCoordinates(string) (Coordinates, error)
}

type Coordinates struct {
	City string  `json:"city"`
	Lat  float64 `json:"latitude"`
	Lon  float64 `json:"longitude"`
}
