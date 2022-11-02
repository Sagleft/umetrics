package main

import (
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
	ChannelContacts *cronContainer
	FindChannels    *cronContainer
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
		color.Red("get channels: %w", err)
		return
	}

	for _, channel := range channels {
		if err := b.saveChannelIFNotExists(channel); err != nil {
			color.Red("save channel: %w", err)
			return
		}
	}
}
