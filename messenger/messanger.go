package messenger

import (
	utopiago "github.com/Sagleft/utopialib-go"
)

type ChannelStats struct{}

type Messenger interface {
	GetStats(channelID string) ChannelStats
}

func NewUtopiaMessenger(clientData utopiago.UtopiaClient) Messenger {
	return &utopia{
		client: &clientData,
	}
}

type utopia struct {
	client *utopiago.UtopiaClient
}

func (u *utopia) GetStats(channelID string) ChannelStats {
	return ChannelStats{}
}
