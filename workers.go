package main

import "log"

func (b *bot) handleJoinChannelTask(event interface{}) {
	e := event.(joinChannelTask)
	if err := b.Messenger.JoinChannel(e.ChannelID, ""); err != nil {
		log.Println(err)
	}
}
