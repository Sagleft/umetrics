package main

import (
	"bot/pkg/memory"
	"log"
	"time"

	"github.com/fatih/color"
	simplecron "github.com/sagleft/simple-cron"
)

type cronContainer struct {
	Cron      *simplecron.CronObject
	InProcess bool
}

func (c *cronContainer) markProcessing(isProcessing bool) {
	c.InProcess = isProcessing
}

type botCrons struct {
	ChannelContacts    *cronContainer
	FindChannels       *cronContainer
	RemoveOldRelations *cronContainer
	FindPeers          *cronContainer
	RemoveOldPeers     *cronContainer
}

func setupCronHandler(callback func(), timeout time.Duration, startImmediate bool) *cronContainer {
	c := &cronContainer{
		Cron: simplecron.NewCronHandler(callback, timeout),
	}
	go c.Cron.Run(startImmediate)
	return c
}

func (b *bot) checkStats() {
	channels, err := b.Memory.GetChannels()
	if err != nil {
		log.Println(err)
		return
	}

	for _, channel := range channels {
		b.Workers.CheckStats.AddEvent(checkChannelTask{
			Channel: channel,
		})
	}
}

func (b *bot) findChannels() {
	channels, err := b.Messenger.GetChannels()
	if err != nil {
		color.Red("get channels: %s", err.Error())
		return
	}

	for _, channel := range channels {
		if err := b.saveChannelIFNotExists(channel); err != nil {
			color.Red("save channel: %s", err.Error())
			return
		}
	}
}

func (b *bot) removeOldRelations() {
	relations, err := b.Memory.GetRelations()
	if err != nil {
		color.Red("get relations: %s", err.Error())
		return
	}

	for _, rel := range relations {
		if time.Since(rel.LastSeen) > maxRelationDuration {
			if err := b.Memory.DeleteRelation(memory.ChannelUserRelation{
				ChannelID:      rel.ChannelID,
				UserPubkeyHash: rel.UserPubkeyHash,
			}); err != nil {
				color.Red("delete relation: %s", err.Error())
				return
			}
		}
	}
}

func (b *bot) findPeers() {

}

func (b *bot) removeOldPeers() {

}
