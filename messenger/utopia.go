package messenger

import (
	"time"

	utopiago "github.com/Sagleft/utopialib-go"
)

func NewUtopiaMessenger(clientData utopiago.UtopiaClient) Messenger {
	return &utopia{
		client: &clientData,
	}
}

type utopia struct {
	client *utopiago.UtopiaClient
}

func (u *utopia) GetStats(channelID string) (ChannelData, error) {
	channelContacts, err := u.client.GetChannelContacts(channelID)
	if err != nil {
		return ChannelData{}, err
	}

	return ChannelData{
		OnlineCount: len(channelContacts),
	}, nil
}

func parseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}

func (u *utopia) GetChannels() ([]SearchChannelData, error) {
	channels, err := u.client.GetChannels(utopiago.GetChannelsTask{})
	if err != nil {
		return nil, err
	}

	r := make([]SearchChannelData, len(channels))
	for i, data := range channels {
		r[i] = SearchChannelData{
			Name:        data.Name,
			ChannelID:   data.ChannelID,
			OwnerPubkey: data.OwnerPubkey,
			IsPrivate:   data.IsPrivate,
			Description: data.Description,
		}

		r[i].CreatedOn, err = parseTime(data.CreatedOn)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}
