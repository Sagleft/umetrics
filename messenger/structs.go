package messenger

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
