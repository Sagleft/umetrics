package memory

type Memory interface {
	IsChannelExists(channelID string) (bool, error)
	SaveChannel(Channel) error
	GetChannelsCount() (int64, error)

	IsUserExists(User) (bool, error)
	SaveUser(User) error

	GetChannels() ([]Channel, error)

	IsRelationExists(ChannelUserRelation) (bool, error)
	SaveRelation(ChannelUserRelation) error
	GetRelations() ([]ChannelUserRelation, error)
	DeleteRelation(ChannelUserRelation) error

	IsPeerExists(Peer) (bool, error)
	SavePeer(Peer) error
	GetPeer(Peer) (Peer, error)
	GetPeers() ([]Peer, error)
	DeletePeer(Peer) error
}
