package main

import (
	"fmt"
	"log"
	"time"

	"bot/config"
	"bot/memory"
	"bot/messenger"

	swissknife "github.com/Sagleft/swiss-knife"
	"github.com/beefsack/go-rate"
	simplecron "github.com/sagleft/simple-cron"
)

const (
	configJSONPath           = "config.json"
	dbFilename               = "memory.db"
	checkChannelsTimeout     = time.Minute * 5
	checkChannelsInStart     = true
	queueDefaultMaxCapacity  = 1000
	limitMaxJoinChannelTasks = 3 * time.Second // per second
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

type queueWorker struct {
	W       *swissknife.ChannelWorker
	Limiter *rate.RateLimiter
}

type queueWorkers struct {
	JoinChannel queueWorker
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
	// setup queues
	b.Workers = queueWorkers{
		JoinChannel: queueWorker{
			W:       swissknife.NewChannelWorker(b.handleJoinChannelTask, queueDefaultMaxCapacity),
			Limiter: rate.New(1, limitMaxJoinChannelTasks),
		},
	}
	go b.Workers.JoinChannel.W.Start()

	// setup cron
	b.Handlers = botCrons{
		Channels: cronContainer{
			Cron: simplecron.NewCronHandler(b.checkChannels, checkChannelsTimeout),
		},
	}
	go b.Handlers.Channels.Cron.Run(checkChannelsInStart)

	// TODO: setup channels online cron
	return nil
}
