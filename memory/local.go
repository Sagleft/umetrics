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

func (db *localDB) isEntryExists(entryPointer interface{}) (bool, error) {
	result := db.conn.First(entryPointer)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}

func (db *localDB) IsChannelExists(channelID string) (bool, error) {
	return db.isEntryExists(&Channel{
		ID: channelID,
	})
}

func (db *localDB) SaveChannel(c Channel) error {
	return db.conn.Save(&c).Error
}

func (db *localDB) IsUserExists(userPubkeyHash string) (bool, error) {
	return db.isEntryExists(&User{
		PubkeyHash: userPubkeyHash,
	})
}

func (db *localDB) SaveUser(u User) error {
	return db.conn.Save(&u).Error
}
