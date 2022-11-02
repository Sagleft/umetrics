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
	checkChannelContactsTimeout = time.Minute * 10
	checkStatsAtStart           = true
	queueDefaultMaxCapacity     = 3000
)

type bot struct {
	Memory        memory.Memory
	Messenger     messenger.Messenger
	Handlers      botCrons
	Workers       queueWorkers
	BotPubkeyHash string
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

	if err := b.loadOwnContact(); err != nil {
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

func (b *bot) loadOwnContact() error {
	ownContact, err := b.Messenger.GetOwnContact()
	if err != nil {
		return err
	}

	b.BotPubkeyHash = ownContact.PubkeyHash
	return nil
}

func (b *bot) run() error {
	// setup queues
	b.Workers = queueWorkers{
		CheckStats: swissknife.NewChannelWorker(
			b.checkStats,
			queueDefaultMaxCapacity,
		).SetAsync(false),
	}
	go b.Workers.CheckStats.Start()

	// setup cron
	b.Handlers = botCrons{
		ChannelContacts: setupCronHandler(b.checkUsers, checkChannelContactsTimeout, checkStatsAtStart),
	}
	return nil
}
