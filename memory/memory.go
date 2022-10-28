package memory

import (
	"fmt"

	swissknife "github.com/Sagleft/swiss-knife"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Memory interface{}

type defaultDB struct {
	conn *gorm.DB
}

func NewLocalDB(filename string) (Memory, error) {
	if !swissknife.IsFileExists(filename) {
		return nil, fmt.Errorf("db file not found: %q", filename)
	}

	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	for _, prefab := range models {
		if err := db.AutoMigrate(prefab); err != nil {
			return nil, fmt.Errorf("failed to migrate: %w", err)
		}
	}

	return &defaultDB{
		conn: db,
	}, nil
}
