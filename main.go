package main

import (
	"bot/memory"
	"bot/messenger"
	"fmt"
	"log"
	"time"

	swissknife "github.com/Sagleft/swiss-knife"
	utopiago "github.com/Sagleft/utopialib-go"
	simplecron "github.com/sagleft/simple-cron"
)

const (
	dbFilename           = "memory.db"
	checkChannelsTimeout = time.Minute * 5
)

func main() {
	b, err := newBot()
	if err != nil {
		log.Fatalln(err)
	}

	if b.run(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Bot started")
	swissknife.RunInBackground()
}

type bot struct {
	Memory    memory.Memory
	Messenger messenger.Messenger
}

func newBot() (*bot, error) {
	db, err := memory.NewLocalDB(dbFilename)
	if err != nil {
		return nil, err
	}

	return &bot{
		Memory:    db,
		Messenger: messenger.NewUtopiaMessenger(utopiago.UtopiaClient{}),
	}, nil
}

func (b *bot) run() error {
	// setup channels list cron
	simplecron.NewCronHandler(b.checkChannels, checkChannelsTimeout).Run(true)

	// TODO: setup channels online cron
	return nil
}

func (b *bot) checkChannels() {
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
			return
		}

		if err := b.Memory.SaveChannel(data); err != nil {
			log.Println(err)
			return
		}
	}
}
