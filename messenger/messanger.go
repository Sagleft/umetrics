package messenger

import (
	"bot/memory"
)

type Messenger interface {
	GetStats(channelID string) (ChannelData, error)
	GetChannels() ([]memory.Channel, error)
}
