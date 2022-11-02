package main

import (
	"bot/memory"
	"fmt"
	"log"
	"time"

	swissknife "github.com/Sagleft/swiss-knife"
	"github.com/beefsack/go-rate"
	"github.com/fatih/color"
)

type queueWorker struct {
	W       *swissknife.ChannelWorker
	Limiter *rate.RateLimiter
}

type queueWorkers struct {
	JoinChannel         queueWorker
	CheckChannelContact queueWorker
}

type joinChannelTask struct {
	ChannelID           string
	EnableNotifications bool
}

type checkChannelTask struct {
	Channel memory.Channel
}

func (b *bot) handleJoinChannelTask(event interface{}) {
	if b.Handlers.Channels.InProcess {
		return
	}
	b.Handlers.Channels.markProcessing(true)
	defer b.Handlers.Channels.markProcessing(false)

	e := event.(joinChannelTask)
	log.Println("join to channel " + e.ChannelID)

	if err := b.Messenger.JoinChannel(e.ChannelID, ""); err != nil {
		log.Println(err)
		return
	}

	if !e.EnableNotifications {
		if err := b.Messenger.ToogleChannelNotifications(e.ChannelID, false); err != nil {
			log.Println(err)
		}
	}
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

	log.Printf("new user saved: %s", u.Nickname)
	fmt.Println()
	return nil
}

func (b *bot) checkChannelContact(event interface{}) {
	if b.Handlers.ChannelContacts.InProcess {
		return
	}
	b.Handlers.ChannelContacts.markProcessing(true)
	defer b.Handlers.ChannelContacts.markProcessing(false)

	e := event.(checkChannelTask)
	log.Println("check channel " + e.Channel.Title + "..")

	if err := b.Messenger.JoinChannel(e.Channel.ID, ""); err != nil {
		color.Red("failed to join to %s: %w", e.Channel.ID, err)
		return
	}

	queryTimestamp := time.Now()
	contacts, err := b.Messenger.GetChannelContacts(e.Channel.ID)
	if err != nil {
		log.Println(err)
		return
	}

	if len(contacts) == 0 {
		return
	}

	for _, contact := range contacts {
		if err := b.saveUserIfNotExists(memory.User{
			PubkeyHash: contact.PubkeyHash,
			Nickname:   contact.Nick,
			LastSeen:   queryTimestamp,
		}); err != nil {
			log.Println(err)
			return
		}
	}
}
