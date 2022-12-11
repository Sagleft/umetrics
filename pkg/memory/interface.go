package memory

import "time"

type Memory interface {
	IsChannelExists(channelID string) (bool, error)
	SaveChannel(Channel) error
	GetChannelsCount() (int64, error)
	GetChannels() ([]Channel, error)
	GetTopChannels(count int) ([]ChannelOnline, error)

	IsUserExists(User) (bool, error)
	AddUser(User) error
	UpdateUserLastSeen(u User, lastSeen time.Time) error
	GetUsersCount() (int64, error)

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
