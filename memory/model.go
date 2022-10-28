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

	PubkeyHash  string `gorm:"type:varchar(32);column:pubkey_hash;unique_index" json:"pubkey_hash"`
	Nickname    string `gorm:"type:varchar(48);default:'anonymous';column:nickname" json:"nickname"`
	IsModerator bool   `gorm:"type:bool;default:false;column:is_moderator" json:"is_moderator"`
}

func (User) TableName() string {
	return "users"
}

/*
original data:
	AvatarID    string `json:"avatarId"`    // example: defAvatar_F10383EA72AC6263C21F356CD8D2E2A2
	ChannelID   string `json:"channelid"`   // F10383EA72AC6263C21F356CD8D2E2A2
	CreatedOn   string `json:"created"`     // 2022-01-28T16:11:39.144Z
	Description string `json:"description"` // can be empty
	IsJoined    bool   `json:"isjoined"`    // false
	IsPrivate   bool   `json:"isprivate"`   // true
	EditedOn    string `json:"modified"`    // 2022-01-28T16:11:39.145Z
	Name        string `json:"name"`        // Monica
	OwnerPubkey string `json:"owner"`       // 1B742E8D8DAE682ADD2568BE25B23F35BA7A8BFC1D5D3BCA0EE219A754A48201

RFC3339     = "2006-01-02T15:04:05Z07:00"
RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"

*/
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
