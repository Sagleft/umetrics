package memory

type Memory interface {
	IsChannelExists(channelID string) (bool, error)
	SaveChannel(Channel) error

	IsUserExists(User) (bool, error)
	SaveUser(User) error

	GetChannels() ([]Channel, error)

	IsRelationExists(ChannelUserRelation) (bool, error)
	SaveRelation(ChannelUserRelation) error
	GetRelations() ([]ChannelUserRelation, error)
}
