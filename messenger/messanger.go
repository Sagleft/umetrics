package messenger

import (
	"bot/memory"
)

type Messenger interface {
	GetStats(channelID string) (ChannelData, error)
	GetChannels() ([]memory.Channel, error)
	JoinChannel(channelID, password string) error
	GetJoinedChannels() (map[string]struct{}, error)
	ToogleChannelNotifications(channelID string, enabled bool) error
}
