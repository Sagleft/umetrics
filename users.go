package main

import "log"

func (b *bot) checkUsers() {
	channels, err := b.Memory.GetChannels()
	if err != nil {
		log.Println(err)
		return
	}

	for _, channel := range channels {
		b.Workers.CheckChannelContact.AddEvent(checkChannelTask{
			Channel: channel,
		})
	}
}
