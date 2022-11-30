package messenger

import (
	"bot/pkg/memory"

	ustructs "github.com/Sagleft/utopialib-go/v2/pkg/structs"
)

type Messenger interface {
	GetChannels() ([]memory.Channel, error)
	GetChannelContacts(channelID string) ([]ustructs.ChannelContactData, error)
	JoinChannel(channelID, password string) error
	GetJoinedChannels() (map[string]struct{}, error)
	ToogleChannelNotifications(channelID string, enabled bool) error
	GetOwnContact() (ustructs.OwnContactData, error)
	GetChannelData(channelID string) (ustructs.ChannelData, error)
	GetNetworkConnections() ([]ustructs.PeerInfo, error)
}
