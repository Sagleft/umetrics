package main

import (
	"bot/memory"
	"log"

	swissknife "github.com/Sagleft/swiss-knife"
	"github.com/beefsack/go-rate"
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
	Title string
	ID    string
}

func (b *bot) handleJoinChannelTask(event interface{}) {
	b.Workers.JoinChannel.Limiter.Wait()

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

	return b.Memory.SaveUser(u)
}

func (b *bot) checkChannelContact(event interface{}) {
	b.Workers.CheckChannelContact.Limiter.Wait()

	e := event.(checkChannelTask)
	log.Println("check channel " + e.Title + "..")

	contacts, err := b.Messenger.GetChannelContacts(e.ID)
	if err != nil {
		log.Println(err)
		return
	}

	for _, contact := range contacts {
		if err := b.saveUserIfNotExists(memory.User{
			PubkeyHash: contact.PubkeyHash,
			Nickname:   contact.Nick,
		}); err != nil {
			log.Println(err)
			return
		}
	}
}
