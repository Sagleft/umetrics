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
	UpdateUserData(u User, lastSeen time.Time, nickname string) error
	GetUsersCount(daysInterval int) (int, error)
	GetTopUsers(count int) ([]UserOnline, error)

	IsRelationExists(ChannelUserRelation) (bool, error)
	SaveRelation(ChannelUserRelation) error
	UpdateRelationLastSeen(r ChannelUserRelation, lastSeen time.Time) error
	GetRelations() ([]ChannelUserRelation, error)
	DeleteRelation(ChannelUserRelation) error

	IsPeerExists(Peer) (bool, error)
	SavePeer(Peer) error
	GetPeer(Peer) (Peer, error)
	GetPeers() ([]Peer, error)
	DeletePeer(Peer) error

	SaveChannelStats(ChannelStats) error

	GetChannelOwners(count int) ([]ChannelOwner, error)
}
