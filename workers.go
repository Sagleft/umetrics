package main

import (
	"bot/memory"
	"fmt"
	"time"

	swissknife "github.com/Sagleft/swiss-knife"
	utopiago "github.com/Sagleft/utopialib-go"
	"github.com/fatih/color"
)

type queueWorkers struct {
	CheckStats *swissknife.ChannelWorker
}

type checkChannelTask struct {
	Channel memory.Channel
}

func (b *bot) saveUserIfNotExists(u memory.User) error {
	isUserKnown, err := b.Memory.IsUserExists(u.PubkeyHash)
	if err != nil {
		return err
	}
	if isUserKnown {
		return nil
	}

	if err := b.Memory.SaveUser(u); err != nil {
		return err
	}

	color.Green("new user saved: %s", u.Nickname)
	return nil
}

func (b *bot) saveChannelIfNotExists(channel memory.Channel) error {
	isExists, err := b.Memory.IsChannelExists(channel.ID)
	if err != nil {
		return err
	}

	if !isExists {
		if err := b.Memory.SaveChannel(channel); err != nil {
			return err
		}
	}
	return nil
}

func (b *bot) checkStats(event interface{}) {
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

	contacts, err := b.Messenger.GetChannelContacts(e.Channel.ID)
	if err != nil {
		color.Red(err.Error())
		return
	}

	e.Channel.LastOnline = len(contacts)

	if err := b.saveChannelIfNotExists(e.Channel); err != nil {
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
	contacts []utopiago.ChannelContactData,
) error {
	queryTimestamp := time.Now()
	fmt.Printf("check channel %s.. %v online", channel.Title, len(contacts))
	fmt.Println()

	for _, contact := range contacts {
		if contact.PubkeyHash == b.BotPubkeyHash {
			continue
		}

		if err := b.saveUserIfNotExists(memory.User{
			PubkeyHash: contact.PubkeyHash,
			Nickname:   contact.Nick,
			LastSeen:   queryTimestamp,
		}); err != nil {
			return err
		}
	}
	return nil
}
