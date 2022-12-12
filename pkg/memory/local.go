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
	fmt.Println("connect to db..")
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

func (db *localDB) GetChannelsCount() (int64, error) {
	var channelsCount int64
	result := db.conn.Model(&Channel{}).Count(&channelsCount)
	return channelsCount, result.Error
}

func (db *localDB) SaveChannel(c Channel) error {
	return db.conn.Save(&c).Error
}

func (db *localDB) IsUserExists(u User) (bool, error) {
	return db.isEntryExists(&u, &User{})
}

func (db *localDB) AddUser(u User) error {
	return db.conn.Save(&u).Error
}

func (db *localDB) UpdateUserLastSeen(u User, lastSeen time.Time) error {
	return db.conn.Model(&u).Where("pubkey_hash", u.PubkeyHash).Update("last_seen", lastSeen).Error
}

func (db *localDB) GetUsersCount() (int64, error) {
	var usersCount int64
	result := db.conn.Model(&User{}).Count(&usersCount)
	return usersCount, result.Error
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

func (db *localDB) GetRelations() ([]ChannelUserRelation, error) {
	relations := []ChannelUserRelation{}
	result := db.conn.Find(&relations)
	return relations, result.Error
}

func (db *localDB) DeleteRelation(r ChannelUserRelation) error {
	result := db.conn.Delete(&r)
	return result.Error
}

func (db *localDB) IsPeerExists(p Peer) (bool, error) {
	return db.isEntryExists(&p, &Peer{})
}

func (db *localDB) SavePeer(p Peer) error {
	return db.conn.Save(&p).Error
}

func (db *localDB) GetPeers() ([]Peer, error) {
	peers := []Peer{}
	result := db.conn.Select("lon", "lat", "city").Limit(maxPeersPerRequest).Find(&peers)
	return peers, result.Error
}

func (db *localDB) DeletePeer(p Peer) error {
	result := db.conn.Where("IP", p.IP).Delete(&p)
	return result.Error
}

func (db *localDB) GetPeer(p Peer) (Peer, error) {
	peer := Peer{}
	result := db.conn.Where(&p).First(&peer)
	return peer, result.Error
}

func (db *localDB) GetTopChannels(count int) ([]ChannelOnline, error) {
	data := []ChannelOnline{}

	result := db.conn.Raw("SELECT COUNT(c.id) AS contactsCount,c.title FROM channels c INNER JOIN channel_contacts cc ON cc.channel_id=c.id GROUP BY c.id ORDER BY contactsCount DESC LIMIT ?", count).Scan(&data)

	return data, result.Error
}

func (db *localDB) GetTopUsers(count int) ([]UserOnline, error) {
	data := []UserOnline{}

	result := db.conn.Raw("SELECT COUNT(cc.id) AS channelsCount,u.nickname FROM users u INNER JOIN channel_contacts cc ON cc.contact_pubkey_hash=u.pubkey_hash GROUP BY u.id ORDER BY channelsCount DESC LIMIT ?", count).Scan(&data)

	return data, result.Error
}

func (db *localDB) UpdateRelationLastSeen(r ChannelUserRelation, lastSeen time.Time) error {
	return db.conn.Model(&r).Where("contact_pubkey_hash", r.UserPubkeyHash).Update("last_seen", lastSeen).Error
}

func (db *localDB) SaveChannelStats(s ChannelStats) error {
	return db.conn.Save(&s).Error
}
