package main

import (
	"fmt"
	"log"
	"time"

	"bot/pkg/config"
	"bot/pkg/frontend"
	"bot/pkg/locator"
	"bot/pkg/memory"
	"bot/pkg/messenger"

	swissknife "github.com/Sagleft/swiss-knife"
)

const (
	configJSONPath = "config.json"
	dbFilename     = "memory.db"

	checkChannelContactsTimeout = time.Minute * 3
	findChannelsTimeout         = time.Minute * 15
	removeOldRelationsTimeout   = time.Minute * 5
	findPeersTimeout            = time.Minute * 5
	removeOldPeersTimeout       = time.Minute * 5

	maxRelationDuration = time.Hour * 24 * 7
	maxPeerDuration     = time.Hour * 24 * 7

	checkStatsAtStart         = true
	findChannelsAtStart       = true
	removeOldRelationsAtStart = false
	findPeersAtStart          = true
	removeOldPeersAtStart     = true

	queueDefaultMaxCapacity = 3000
)

type bot struct {
	Frontend      frontend.Frontend
	Memory        memory.Memory
	Locator       locator.Locator
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

	fmt.Println("create bot..")
	b, err := newBot(cfg, db)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("start..")
	go b.run()

	fmt.Println("setup frontend..")
	b.Frontend, err = frontend.NewGINFrontend(db)
	if err != nil {
		log.Fatalln(err)
	}

	if err := b.Frontend.Run(); err != nil {
		log.Fatalln(err)
	}
}

func newBot(cfg config.Config, db memory.Memory) (*bot, error) {
	return &bot{
		Memory:    db,
		Locator:   locator.NewDefaultLocator(),
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

func (b *bot) run() {
	fmt.Println("load account data..")
	if err := b.loadOwnContact(); err != nil {
		log.Println(err)
		return
	}

	// setup queues
	b.Workers = queueWorkers{
		CheckStats: swissknife.NewChannelWorker(
			b.handleCheckStatsTask,
			queueDefaultMaxCapacity,
		).SetAsync(false),
	}
	go b.Workers.CheckStats.Start()

	// setup cron
	b.Handlers = botCrons{
		ChannelContacts:    setupCronHandler(b.checkStats, checkChannelContactsTimeout, checkStatsAtStart),
		FindChannels:       setupCronHandler(b.findChannels, findChannelsTimeout, findChannelsAtStart),
		RemoveOldRelations: setupCronHandler(b.removeOldRelations, removeOldRelationsTimeout, removeOldRelationsAtStart),
		FindPeers:          setupCronHandler(b.findPeers, findPeersTimeout, findPeersAtStart),
		RemoveOldPeers:     setupCronHandler(b.removeOldPeers, removeOldPeersTimeout, removeOldPeersAtStart),
	}
}
