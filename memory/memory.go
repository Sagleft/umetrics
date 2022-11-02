package memory

type Memory interface {
	IsChannelExists(channelID string) (bool, error)
	SaveChannel(Channel) error

	IsUserExists(userPubkeyHash string) (bool, error)
	SaveUser(User) error
}
