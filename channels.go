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

	joinedChannels, err := b.Messenger.GetJoinedChannels()
	if err != nil {
		log.Println(err)
		return
	}

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

		if _, isJoined := joinedChannels[data.ID]; !isJoined {
			b.addJoinChannelTask(joinChannelTask{ChannelID: data.ID})
		}
	}
}

type joinChannelTask struct {
	ChannelID           string
	EnableNotifications bool
}

func (b *bot) addJoinChannelTask(task joinChannelTask) {
	log.Println("add joinchannel task")

	b.Workers.JoinChannel.W.AddEvent(task)
}
