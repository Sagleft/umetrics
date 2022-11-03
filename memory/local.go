package memory

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type localDB struct {
	conn *gorm.DB
}

func NewLocalDB(filename string) (Memory, error) {
	lg := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{
		Logger: lg,
	})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	fmt.Println("migrate..")
	for _, prefab := range models {
		if err := db.AutoMigrate(prefab); err != nil {
			return nil, fmt.Errorf("failed to migrate: %w", err)
		}
	}

	return &localDB{
		conn: db,
	}, nil
}

func (db *localDB) isEntryExists(entryPointer interface{}, typePointer interface{}) (bool, error) {
	result := db.conn.Where(entryPointer).First(typePointer)
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
	}, &Channel{})
}

func (db *localDB) SaveChannel(c Channel) error {
	return db.conn.Save(&c).Error
}

func (db *localDB) IsUserExists(u User) (bool, error) {
	return db.isEntryExists(&u, &User{})
}

func (db *localDB) SaveUser(u User) error {
	return db.conn.Save(&u).Error
}

func (db *localDB) GetChannels() ([]Channel, error) {
	channels := []Channel{}
	result := db.conn.Order("last_online desc").Find(&channels)
	return channels, result.Error
}

func (db *localDB) SaveRelation(c ChannelUserRelation) error {
	return db.conn.Save(&c).Error
}

func (db *localDB) IsRelationExists(r ChannelUserRelation) (bool, error) {
	return db.isEntryExists(&r, &ChannelUserRelation{})
}
