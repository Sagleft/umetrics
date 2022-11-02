package memory

import (
	"time"

	"gorm.io/gorm"
)

var models = []interface{}{
	&User{}, &Channel{},
}

type User struct {
	gorm.Model

	PubkeyHash string    `gorm:"type:varchar(32);column:pubkey_hash;unique_index" json:"pubkey_hash"`
	Nickname   string    `gorm:"type:varchar(48);default:'anonymous';column:nickname" json:"nickname"`
	LastSeen   time.Time `gorm:"column:last_seen" json:"last_seen"`
}

func (User) TableName() string {
	return "users"
}

type Channel struct {
	gorm.Model

	ID              string    `gorm:"type:varchar(32);column:id;unique_index"`
	Title           string    `gorm:"type:varchar(96);column:title"`
	Description     string    `gorm:"type:varchar(256);column:description"`
	OwnerPubkey     string    `gorm:"type:varchar(64);column:owner_pubkey"`
	OwnerPubkeyHash string    `gorm:"type:varchar(32);column:owner_hash"`
	IsPrivate       bool      `gorm:"type:bool;default:false;column:is_private"`
	CreatedOn       time.Time `gorm:"column:created_on"`
	LastOnline      int       `gorm:"column:last_online;default:0"`
	GeoTag          UGeoTag   `gorm:"column:geotag;default:''"`
	ReadOnly        bool      `gorm:"column:readonly;default:false"`
	ReadOnlyPrivacy bool      `gorm:"column:readonly_privacy;default:false"`
}

func (Channel) TableName() string {
	return "channels"
}

type ChannelContact struct {
	gorm.Model

	UserPubkeyHash string `gorm:"type:varchar(32);column:contact_pubkey_hash" json:"contact_pubkey_hash"`
	IsModerator    bool   `gorm:"type:bool;default:false;column:is_moderator" json:"is_moderator"`
}

func (ChannelContact) TableName() string {
	return "channel_contacts"
}
