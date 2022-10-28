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
	checkChannelsInStart = false
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

type cronContainer struct {
	Cron      *simplecron.CronObject
	InProcess bool
}

func (c *cronContainer) markProcessing(isProcessing bool) {
	c.InProcess = isProcessing
}

type bot struct {
	Memory    memory.Memory
	Messenger messenger.Messenger

	ChannelsCron cronContainer
}

func newBot(cfg config.Config, db memory.Memory) (*bot, error) {
	return &bot{
		Memory:    db,
		Messenger: messenger.NewUtopiaMessenger(cfg.Utopia),
	}, nil
}

func (b *bot) run() error {
	// setup channels list cron
	b.ChannelsCron = cronContainer{
		Cron: simplecron.NewCronHandler(b.checkChannels, checkChannelsTimeout),
	}
	b.ChannelsCron.Cron.Run(checkChannelsInStart)

	// TODO: setup channels online cron
	return nil
}

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
