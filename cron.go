package main

import (
	"time"

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
}

func setupCronHandler(callback func(), timeout time.Duration, startImmediate bool) *cronContainer {
	c := &cronContainer{
		Cron: simplecron.NewCronHandler(callback, timeout),
	}
	go c.Cron.Run(startImmediate)
	return c
}
