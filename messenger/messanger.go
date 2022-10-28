package messenger

type Messenger interface {
	GetStats(channelID string) (ChannelData, error)
	GetChannels() ([]SearchChannelData, error)
}
