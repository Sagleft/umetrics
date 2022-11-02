package main

import (
	"fmt"
	"log"
)

func (b *bot) checkChannels() {
	if b.Handlers.Channels.InProcess {
		return
	}
	b.Handlers.Channels.markProcessing(true)
	defer b.Handlers.Channels.markProcessing(false)

	channels, err := b.Messenger.GetChannels()
	if err != nil {
		log.Println(err)
	}

	joinedChannels, err := b.Messenger.GetJoinedChannels()
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("joined to %v channels", len(joinedChannels))
	fmt.Println()

	for _, data := range channels {
		isExists, err := b.Memory.IsChannelExists(data.ID)
		if err != nil {
			log.Println(err)
			return
		}

		if !isExists {
			if err := b.Memory.SaveChannel(data); err != nil {
				log.Println(err)
				return
			}
		}

		if !data.IsPrivate {
			if _, isJoined := joinedChannels[data.ID]; !isJoined {
				b.addJoinChannelTask(joinChannelTask{ChannelID: data.ID})
			}
		}
	}
}

func (b *bot) addJoinChannelTask(task joinChannelTask) {
	b.Workers.JoinChannel.AddEvent(task)
}
