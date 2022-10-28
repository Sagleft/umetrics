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
	"gorm.io/gorm"
)

const (
	dbFilename           = "memory.db"
	checkChannelsTimeout = time.Minute * 5
)

func main() {
	db, err := memory.InitDB(dbFilename)
	if err != nil {
		log.Fatalln(err)
	}

	if err := newBot(db).run(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Bot started")
	swissknife.RunInBackground()
}

type bot struct {
	DB        *gorm.DB
	Messenger messenger.Messenger
}

func newBot(db *gorm.DB) *bot {
	return &bot{
		DB:        db,
		Messenger: messenger.NewUtopiaMessenger(utopiago.UtopiaClient{}),
	}
}

func (b *bot) run() error {
	simplecron.NewCronHandler(b.checkChannels, checkChannelsTimeout).Run(true)

	// TODO
	return nil
}

func (b *bot) checkChannels() {

}
