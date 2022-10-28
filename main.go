package main

import (
	"bot/messenger"
	"log"

	utopiago "github.com/Sagleft/utopialib-go"
)

func main() {
	if err := newBot().run(); err != nil {
		log.Fatalln(err)
	}
}

type bot struct {
	Messenger messenger.Messenger
}

func newBot() *bot {
	return &bot{
		Messenger: messenger.NewUtopiaMessenger(utopiago.UtopiaClient{}),
	}
}

func (b *bot) run() error {
	// TODO
	return nil
}
