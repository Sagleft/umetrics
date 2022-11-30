package messenger

import (
	"bot/pkg/memory"
	"time"

	utopiago "github.com/Sagleft/utopialib-go"
	"github.com/beefsack/go-rate"
)

const (
	maxGetChannelsTasks        = 5
	maxGetChannelDataTasks     = 5
	maxJoinChannelTasks        = 5  // per second
	maxGetChannelContactsTasks = 10 // per second
)

type utopia struct {
	client   *utopiago.UtopiaClient
	limiters rateLimiters
}

type rateLimiters struct {
	GetChannels        *rate.RateLimiter
	JoinChannel        *rate.RateLimiter
	GetChannelContacts *rate.RateLimiter
	GetChannelData     *rate.RateLimiter
}

func NewUtopiaMessenger(clientData utopiago.UtopiaClient) Messenger {
	return &utopia{
		client: &clientData,
		limiters: rateLimiters{
			GetChannels:        rate.New(maxGetChannelsTasks, time.Second),
			JoinChannel:        rate.New(maxJoinChannelTasks, time.Second),
			GetChannelContacts: rate.New(maxGetChannelContactsTasks, time.Second),
			GetChannelData:     rate.New(maxGetChannelDataTasks, time.Second),
		},
	}
}

func (u *utopia) GetChannels() ([]memory.Channel, error) {
	u.limiters.GetChannels.Wait()

	channels, err := u.client.GetChannels(utopiago.GetChannelsTask{
		SortBy: utopiago.SortChannelsByModified,
	})
	if err == nil {
		return convertChannelsData(channels)
	}

	if !utopiago.CheckErrorConnBroken(err) {
		return nil, err
	}

	reconnect(func() error {
		channels, err = u.client.GetChannels(utopiago.GetChannelsTask{
			SortBy: utopiago.SortChannelsByModified,
		})
		return err
	})

	return convertChannelsData(channels)
}

func convertChannelsData(channels []utopiago.SearchChannelData) ([]memory.Channel, error) {
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

func (u *utopia) GetChannelContacts(channelID string) ([]utopiago.ChannelContactData, error) {
	u.limiters.GetChannelContacts.Wait()

	contacts, err := u.client.GetChannelContacts(channelID)
	if err == nil {
		return contacts, nil
	}

	if !utopiago.CheckErrorConnBroken(err) {
		return nil, err
	}

	reconnect(func() error {
		contacts, err = u.client.GetChannelContacts(channelID)
		return err
	})
	return contacts, nil
}

func (u *utopia) JoinChannel(channelID, password string) error {
	u.limiters.JoinChannel.Wait()

	_, err := u.client.JoinChannel(channelID, password)
	if err == nil {
		return nil
	}

	if !utopiago.CheckErrorConnBroken(err) {
		return err
	}

	reconnect(func() error {
		_, err := u.client.JoinChannel(channelID, password)
		return err
	})
	return nil
}

func (u *utopia) GetJoinedChannels() (map[string]struct{}, error) {
	channels, err := u.client.GetChannels(utopiago.GetChannelsTask{
		ChannelType: utopiago.ChannelTypeJoined,
	})
	if err != nil {
		if !utopiago.CheckErrorConnBroken(err) {
			return nil, err
		}

		reconnect(func() error {
			channels, err = u.client.GetChannels(utopiago.GetChannelsTask{
				ChannelType: utopiago.ChannelTypeJoined,
			})
			return err
		})
		return convertChannelIDs(channels), nil
	}

	return convertChannelIDs(channels), nil
}

func convertChannelIDs(channels []utopiago.SearchChannelData) map[string]struct{} {
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

	if !utopiago.CheckErrorConnBroken(err) {
		return err
	}

	reconnect(func() error {
		return u.client.ToogleChannelNotifications(channelID, enabled)
	})
	return nil
}

func (u *utopia) GetOwnContact() (utopiago.OwnContactData, error) {
	data, err := u.client.GetOwnContact()
	if err != nil {
		if !utopiago.CheckErrorConnBroken(err) {
			return utopiago.OwnContactData{}, err
		}

		reconnect(func() error {
			data, err = u.client.GetOwnContact()
			return err
		})
		return data, nil
	}

	return data, nil
}

func (u *utopia) GetChannelData(channelID string) (utopiago.ChannelData, error) {
	u.limiters.GetChannelData.Wait()

	data, err := u.client.GetChannelInfo(channelID)
	if err != nil {
		if !utopiago.CheckErrorConnBroken(err) {
			return utopiago.ChannelData{}, err
		}

		reconnect(func() error {
			data, err = u.client.GetChannelInfo(channelID)
			return err
		})
		return data, nil
	}

	return data, nil
}
