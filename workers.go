package main

import "log"

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
