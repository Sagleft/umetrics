package messenger

import (
	"bot/memory"

	utopiago "github.com/Sagleft/utopialib-go"
)

type Messenger interface {
	GetChannels(channelID string) ([]memory.Channel, error)
	GetChannelContacts() ([]utopiago.ChannelContactData, error)
	JoinChannel(channelID, password string) error
	GetJoinedChannels() (map[string]struct{}, error)
	ToogleChannelNotifications(channelID string, enabled bool) error
}
