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
	configJSONPath          = "config.json"
	dbFilename              = "memory.db"
	checkChannelsTimeout    = time.Minute * 5
	checkChannelsInStart    = false
	queueDefaultMaxCapacity = 1000
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
	Handlers  botCrons
	Workers   queueWorkers
}

type queueWorkers struct {
	JoinChannel *swissknife.ChannelWorker
}

type botCrons struct {
	Channels cronContainer
}

func newBot(cfg config.Config, db memory.Memory) (*bot, error) {
	return &bot{
		Memory:    db,
		Messenger: messenger.NewUtopiaMessenger(cfg.Utopia),
	}, nil
}

func (b *bot) run() error {
	b.Workers.JoinChannel = swissknife.NewChannelWorker(b.handleJoinChannelTask, queueDefaultMaxCapacity)
	go b.Workers.JoinChannel.Start()

	b.Handlers = botCrons{
		Channels: cronContainer{
			Cron: simplecron.NewCronHandler(b.checkChannels, checkChannelsTimeout),
		},
	}

	go b.Handlers.Channels.Cron.Run(checkChannelsInStart)

	// TODO: setup channels online cron
	return nil
}
