package messenger

import "time"

/*
original data:
	PubkeyHash  string `json:"hashedPk"`
	LastSeen    string `json:"lastSeen"`
	IsLocal     bool   `json:"local"`
	IsModerator bool   `json:"moderator"`
	Nick        string `json:"nick"`
	Pubkey      string `json:"pk"`
*/
type ChannelContactData struct {
	PubkeyHash  string
	IsModerator bool
	Nick        string
}

type ChannelData struct {
	OnlineCount int
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
type SearchChannelData struct {
	Title           string
	ChannelID       string
	OwnerPubkey     string
	OwnerPubkeyHash string
	IsPrivate       bool
	Description     string
	CreatedOn       time.Time
}
