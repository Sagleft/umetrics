package config

import (
	utopiago "github.com/Sagleft/utopialib-go"
)

type Config struct {
	Utopia utopiago.UtopiaClient `json:"utopia"`
}
