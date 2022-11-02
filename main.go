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
)

const (
	configJSONPath              = "config.json"
	dbFilename                  = "memory.db"
	checkChannelsTimeout        = time.Minute * 10
	checkChannelContactsTimeout = time.Minute * 5
	checkChannelsInStart        = false
	checkContactsInStart        = true
	queueDefaultMaxCapacity     = 1000
	limitMaxJoinChannelTasks    = 3  // per second
	limitMaxCheckChannelTasks   = 30 // per second
)

type bot struct {
	Memory    memory.Memory
	Messenger messenger.Messenger
	Handlers  botCrons
	Workers   queueWorkers
}

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
			W:       swissknife.NewChannelWorker(b.handleJoinChannelTask, queueDefaultMaxCapacity).SetAsync(false),
			Limiter: rate.New(limitMaxJoinChannelTasks, time.Second),
		},
		CheckChannelContact: queueWorker{
			W:       swissknife.NewChannelWorker(b.checkChannelContact, queueDefaultMaxCapacity).SetAsync(false),
			Limiter: rate.New(limitMaxCheckChannelTasks, time.Second),
		},
	}
	go b.Workers.JoinChannel.W.Start()
	go b.Workers.CheckChannelContact.W.Start()

	// setup cron
	b.Handlers = botCrons{
		Channels:        setupCronHandler(b.checkChannels, checkChannelsTimeout, checkChannelsInStart),
		ChannelContacts: setupCronHandler(b.checkUsers, checkChannelContactsTimeout, checkContactsInStart),
	}
	return nil
}
