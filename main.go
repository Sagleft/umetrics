package main

import (
	"fmt"
	"log"
	"time"

	"bot/config"
	"bot/memory"
	"bot/messenger"

	swissknife "github.com/Sagleft/swiss-knife"
)

const (
	configJSONPath              = "config.json"
	dbFilename                  = "memory.db"
	checkChannelsTimeout        = time.Minute * 10
	checkChannelContactsTimeout = time.Minute * 5
	checkChannelsInStart        = false
	checkContactsInStart        = true
	queueDefaultMaxCapacity     = 1000
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
	go swissknife.NewChannelWorker(b.handleJoinChannelTask, queueDefaultMaxCapacity).SetAsync(false).Start()
	go swissknife.NewChannelWorker(b.checkChannelContact, queueDefaultMaxCapacity).SetAsync(false)

	/*b.Workers = queueWorkers{
		JoinChannel:         swissknife.NewChannelWorker(b.handleJoinChannelTask, queueDefaultMaxCapacity).SetAsync(false),
		CheckChannelContact: swissknife.NewChannelWorker(b.checkChannelContact, queueDefaultMaxCapacity).SetAsync(false),
	}
	go b.Workers.JoinChannel.Start()
	go b.Workers.CheckChannelContact.Start()*/

	// setup cron
	b.Handlers = botCrons{
		Channels:        setupCronHandler(b.checkChannels, checkChannelsTimeout, checkChannelsInStart),
		ChannelContacts: setupCronHandler(b.checkUsers, checkChannelContactsTimeout, checkContactsInStart),
	}
	return nil
}
