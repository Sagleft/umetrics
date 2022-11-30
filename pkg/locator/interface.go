package locator

type Locator interface {
	GetCoordinates(string) (Coordinates, error)
}

type Coordinates struct {
	Lat float64 `json:"latitude"`
	Lon float64 `json:"longitude"`
}
