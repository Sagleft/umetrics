package main

import (
	"bot/memory"
	"bot/messenger"
	"log"

	utopiago "github.com/Sagleft/utopialib-go"
	"gorm.io/gorm"
)

const (
	dbFilename = "memory.db"
)

func main() {
	db, err := memory.InitDB(dbFilename)
	if err != nil {
		log.Fatalln(err)
	}

	if err := newBot(db).run(); err != nil {
		log.Fatalln(err)
	}
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
	// TODO
	return nil
}
