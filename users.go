package main

import "log"

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
