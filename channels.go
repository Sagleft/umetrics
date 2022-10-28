package main

import "log"

func (b *bot) checkChannels() {
	if b.ChannelsCron.InProcess {
		return
	}
	b.ChannelsCron.markProcessing(true)
	defer b.ChannelsCron.markProcessing(false)

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
