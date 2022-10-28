package memory

import (
	"gorm.io/gorm"
)

var models = []interface{}{
	&User{},
}

type User struct {
	gorm.Model

	PubkeyHash  string
	Nickname    string
	IsModerator bool
}
