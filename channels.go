package main

import "log"

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

	for _, data := range channels {
		isExists, err := b.Memory.IsChannelExists(data.ID)
		if err != nil {
			log.Println(err)
			return
		}

		if isExists {
			continue
		}

		if err := b.Memory.SaveChannel(data); err != nil {
			log.Println(err)
			return
		}
	}
}

type joinChannelTask struct {
	ChannelID string
}

func (b *bot) addJoinChannelTask(task joinChannelTask) {
	b.Workers.JoinChannel.AddEvent(task)
}
