package messenger

import (
	"time"

	"bot/memory"

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

func (u *utopia) GetChannels() ([]memory.Channel, error) {
	channels, err := u.client.GetChannels(utopiago.GetChannelsTask{})
	if err != nil {
		return nil, err
	}

	r := make([]memory.Channel, len(channels))
	for i, data := range channels {
		r[i] = memory.Channel{
			ID:              data.ChannelID,
			Title:           data.Name,
			OwnerPubkey:     data.OwnerPubkey,
			OwnerPubkeyHash: getMD5Hash(data.OwnerPubkey),
			IsPrivate:       data.IsPrivate,
			Description:     data.Description,
		}

		r[i].CreatedOn, err = parseTime(data.CreatedOn)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (u *utopia) JoinChannel(channelID, password string) error {
	var err error
	if password == "" {
		_, err = u.client.JoinChannel(channelID)
	} else {
		_, err = u.client.JoinChannel(channelID, password)
	}
	return err
}

func (u *utopia) GetJoinedChannels() (map[string]struct{}, error) {
	channels, err := u.client.GetChannels(utopiago.GetChannelsTask{
		ChannelType: utopiago.ChannelTypeJoined,
	})
	if err != nil {
		return nil, err
	}

	channelIDs := make(map[string]struct{})
	for _, data := range channels {
		channelIDs[data.ChannelID] = struct{}{}
	}
	return channelIDs, nil
}

func (u *utopia) ToogleChannelNotifications(channelID string, enabled bool) error {
	return u.client.ToogleChannelNotifications(channelID, enabled)
}
