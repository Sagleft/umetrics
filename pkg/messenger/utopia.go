package messenger

import (
	"bot/pkg/memory"

	utopiago "github.com/Sagleft/utopialib-go/v2"
	uconsts "github.com/Sagleft/utopialib-go/v2/pkg/consts"
	uerrors "github.com/Sagleft/utopialib-go/v2/pkg/errors"
	ustructs "github.com/Sagleft/utopialib-go/v2/pkg/structs"
)

type utopia struct {
	client utopiago.Client
}

func NewUtopiaMessenger(clientData utopiago.Config) Messenger {
	return &utopia{
		client: utopiago.NewUtopiaClient(clientData),
	}
}

func (u *utopia) GetChannels() ([]memory.Channel, error) {
	channels, err := u.client.GetChannels(ustructs.GetChannelsTask{
		SortBy: uconsts.SortChannelsByModified,
	})
	if err == nil {
		return convertChannelsData(channels)
	}

	if !uerrors.CheckErrorConnBroken(err) {
		return nil, err
	}

	reconnect(func() error {
		channels, err = u.client.GetChannels(ustructs.GetChannelsTask{
			SortBy: uconsts.SortChannelsByModified,
		})
		return err
	})

	return convertChannelsData(channels)
}

func convertChannelsData(channels []ustructs.SearchChannelData) ([]memory.Channel, error) {
	var err error
	r := make([]memory.Channel, len(channels))
	for i := len(channels) - 1; i >= 0; i-- {
		data := channels[i]
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

func (u *utopia) GetChannelContacts(channelID string) ([]ustructs.ChannelContactData, error) {
	contacts, err := u.client.GetChannelContacts(channelID)
	if err == nil {
		return contacts, nil
	}

	if !uerrors.CheckErrorConnBroken(err) {
		return nil, err
	}

	reconnect(func() error {
		contacts, err = u.client.GetChannelContacts(channelID)
		return err
	})
	return contacts, nil
}

func (u *utopia) JoinChannel(channelID, password string) error {
	_, err := u.client.JoinChannel(channelID, password)
	if err == nil {
		return nil
	}

	if !uerrors.CheckErrorConnBroken(err) {
		return err
	}

	reconnect(func() error {
		_, err := u.client.JoinChannel(channelID, password)
		return err
	})
	return nil
}

func (u *utopia) GetJoinedChannels() (map[string]struct{}, error) {
	channels, err := u.client.GetChannels(ustructs.GetChannelsTask{
		ChannelType: uconsts.ChannelTypeJoined,
	})
	if err != nil {
		if !uerrors.CheckErrorConnBroken(err) {
			return nil, err
		}

		reconnect(func() error {
			channels, err = u.client.GetChannels(ustructs.GetChannelsTask{
				ChannelType: uconsts.ChannelTypeJoined,
			})
			return err
		})
		return convertChannelIDs(channels), nil
	}

	return convertChannelIDs(channels), nil
}

func convertChannelIDs(channels []ustructs.SearchChannelData) map[string]struct{} {
	channelIDs := make(map[string]struct{})
	for _, data := range channels {
		channelIDs[data.ChannelID] = struct{}{}
	}
	return channelIDs
}

func (u *utopia) ToogleChannelNotifications(channelID string, enabled bool) error {
	err := u.client.ToogleChannelNotifications(channelID, enabled)
	if err == nil {
		return nil
	}

	if !uerrors.CheckErrorConnBroken(err) {
		return err
	}

	reconnect(func() error {
		return u.client.ToogleChannelNotifications(channelID, enabled)
	})
	return nil
}

func (u *utopia) GetOwnContact() (ustructs.OwnContactData, error) {
	data, err := u.client.GetOwnContact()
	if err != nil {
		if !uerrors.CheckErrorConnBroken(err) {
			return ustructs.OwnContactData{}, err
		}

		reconnect(func() error {
			data, err = u.client.GetOwnContact()
			return err
		})
		return data, nil
	}

	return data, nil
}

func (u *utopia) GetChannelData(channelID string) (ustructs.ChannelData, error) {
	data, err := u.client.GetChannelInfo(channelID)
	if err != nil {
		if !uerrors.CheckErrorConnBroken(err) {
			return ustructs.ChannelData{}, err
		}

		reconnect(func() error {
			data, err = u.client.GetChannelInfo(channelID)
			return err
		})
		return data, nil
	}

	return data, nil
}

func (u *utopia) GetNetworkConnections() ([]ustructs.PeerInfo, error) {
	data, err := u.client.GetNetworkConnections()
	if err == nil {
		return data, nil
	}

	if !uerrors.CheckErrorConnBroken(err) {
		return nil, err
	}

	reconnect(func() error {
		data, err = u.client.GetNetworkConnections()
		return err
	})
	return data, nil
}
