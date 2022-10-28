package main

import (
	"fmt"
	"log"
	"time"

	"bot/config"
	"bot/memory"
	"bot/messenger"

	swissknife "github.com/Sagleft/swiss-knife"
	simplecron "github.com/sagleft/simple-cron"
)

const (
	configJSONPath       = "config.json"
	dbFilename           = "memory.db"
	checkChannelsTimeout = time.Minute * 5
)

func main() {
	cfg, err := config.Parse(configJSONPath)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := memory.NewLocalDB(dbFilename)
	if err != nil {
		log.Fatalln(err)
	}

	b, err := newBot(cfg, db)
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

func newBot(cfg config.Config, db memory.Memory) (*bot, error) {
	return &bot{
		Memory:    db,
		Messenger: messenger.NewUtopiaMessenger(cfg.Utopia),
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
