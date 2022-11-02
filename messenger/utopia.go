package messenger

import (
	"bot/memory"
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
	if err != nil {
		return nil, err
	}

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

	return u.client.GetChannelContacts(channelID)
}

func (u *utopia) JoinChannel(channelID, password string) error {
	u.limiters.JoinChannel.Wait()

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

func (u *utopia) GetOwnContact() (utopiago.OwnContactData, error) {
	return u.client.GetOwnContact()
}

func (u *utopia) GetChannelData(channelID string) (utopiago.ChannelData, error) {
	u.limiters.GetChannelData.Wait()

	return u.client.GetChannelInfo(channelID)
}
