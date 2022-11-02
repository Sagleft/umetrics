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

	PubkeyHash string `gorm:"type:varchar(32);column:pubkey_hash;unique_index" json:"pubkey_hash"`
	Nickname   string `gorm:"type:varchar(48);default:'anonymous';column:nickname" json:"nickname"`
}

func (User) TableName() string {
	return "users"
}

type Channel struct {
	gorm.Model

	ID              string    `gorm:"type:varchar(32);column:id;unique_index" json:"id"`
	Title           string    `gorm:"type:varchar(96);column:title" json:"title"`
	Description     string    `gorm:"type:varchar(256);column:description" json:"description"`
	OwnerPubkey     string    `gorm:"type:varchar(64);column:owner_pubkey" json:"owner_pubkey"`
	OwnerPubkeyHash string    `gorm:"type:varchar(32);column:owner_hash" json:"owner_hash"`
	IsPrivate       bool      `gorm:"type:bool;default:false;column:is_private" json:"is_private"`
	CreatedOn       time.Time `gorm:"column:created_on" json:"created_on"`
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
