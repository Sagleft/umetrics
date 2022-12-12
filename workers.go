package main

import (
	"bot/pkg/memory"
	"time"

	swissknife "github.com/Sagleft/swiss-knife"
	ustructs "github.com/Sagleft/utopialib-go/v2/pkg/structs"
	"github.com/fatih/color"
)

type queueWorkers struct {
	CheckStats *swissknife.ChannelWorker
}

type checkChannelTask struct {
	Channel memory.Channel
}

func (b *bot) saveUser(u memory.User) error {
	isUserKnown, err := b.Memory.IsUserExists(memory.User{
		PubkeyHash: u.PubkeyHash,
	})
	if err != nil {
		return err
	}
	if isUserKnown {
		return b.Memory.UpdateUserLastSeen(u, time.Now())
	}

	if err := b.Memory.AddUser(u); err != nil {
		return err
	}
	color.Green("new user saved: %s", u.Nickname)
	return nil
}

func (b *bot) saveChannelIFNotExists(channel memory.Channel) error {
	isExists, err := b.Memory.IsChannelExists(channel.ID)
	if err != nil {
		return err
	}
	if isExists {
		return nil
	}

	color.Green("save new channel: %s", channel.Title)
	return b.Memory.SaveChannel(channel)
}

func (b *bot) saveUserRelation(channel memory.Channel, contact ustructs.ChannelContactData) error {
	isExists, err := b.Memory.IsRelationExists(memory.ChannelUserRelation{
		ChannelID:      channel.ID,
		UserPubkeyHash: contact.PubkeyHash,
	})
	if err != nil {
		return err
	}
	if isExists {
		b.Memory.UpdateRelationLastSeen(memory.ChannelUserRelation{
			UserPubkeyHash: contact.PubkeyHash,
		}, time.Now())
		return nil
	}

	if err := b.Memory.SaveRelation(memory.ChannelUserRelation{
		ChannelID:      channel.ID,
		UserPubkeyHash: contact.PubkeyHash,
		IsModerator:    contact.IsModerator,
		LastSeen:       time.Now(),
	}); err != nil {
		return err
	}

	if !isExists {
		color.Green("new relation saved: %s is a member of %q", contact.Nick, channel.Title)
	}
	return nil
}

func (b *bot) handleCheckStatsTask(event interface{}) {
	if b.Handlers.ChannelContacts.InProcess {
		return
	}
	b.Handlers.ChannelContacts.markProcessing(true)
	defer b.Handlers.ChannelContacts.markProcessing(false)

	e := event.(checkChannelTask)
	if e.Channel.IsPrivate {
		return // ignore private channels
	}

	if err := b.Messenger.JoinChannel(e.Channel.ID, ""); err != nil {
		color.Red(err.Error())
		return
	}

	// get channel online
	contacts, err := b.Messenger.GetChannelContacts(e.Channel.ID)
	if err != nil {
		color.Red(err.Error())
		return
	}
	e.Channel.LastOnline = len(contacts)

	channelData, err := b.Messenger.GetChannelData(e.Channel.ID)
	if err != nil {
		color.Red(err.Error())
		return
	}
	e.Channel.GeoTag = memory.UGeoTag(channelData.GeoTag)
	e.Channel.ReadOnly = channelData.ReadOnly
	e.Channel.ReadOnlyPrivacy = channelData.ReadOnlyPrivacy

	if err := b.Memory.SaveChannel(e.Channel); err != nil {
		color.Red(err.Error())
		return
	}

	// PROCESS CHANNEL CONTACTS
	if len(contacts) > 0 {
		if err := b.processChannelContacts(e.Channel, contacts); err != nil {
			color.Red(err.Error())
			return
		}
	}
}

func (b *bot) processChannelContacts(
	channel memory.Channel,
	contacts []ustructs.ChannelContactData,
) error {
	queryTimestamp := time.Now()
	//fmt.Printf("check channel %s (%s).. %v online", channel.Title, channel.ID, len(contacts))
	//fmt.Println()

	for _, contact := range contacts {
		if contact.PubkeyHash == b.BotPubkeyHash {
			continue
		}

		if err := b.saveUser(memory.User{
			PubkeyHash: contact.PubkeyHash,
			Nickname:   contact.Nick,
			LastSeen:   queryTimestamp,
		}); err != nil {
			return err
		}

		if err := b.saveUserRelation(channel, contact); err != nil {
			return err
		}
	}
	return nil
}
