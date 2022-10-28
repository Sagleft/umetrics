package memory

import (
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type localDB struct {
	conn *gorm.DB
}

func NewLocalDB(filename string) (Memory, error) {
	/*if !swissknife.IsFileExists(filename) {
		return nil, fmt.Errorf("db file not found: %q", filename)
	}*/

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

	return &localDB{
		conn: db,
	}, nil
}

func (db *localDB) IsChannelExists(channelID string) (bool, error) {
	channelEntry := Channel{
		ID: channelID,
	}

	result := db.conn.First(&channelEntry)
	if result.Error != nil {
		return false, result.Error
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return true, nil
}

func (db *localDB) SaveChannel(c Channel) error {
	return db.conn.Save(&c).Error
}
