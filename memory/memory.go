package memory

type Memory interface {
	IsChannelExists(channelID string) (bool, error)
	SaveChannel(Channel) error
}
