package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// load config
	cfg, err := LoadConfig("config")
	if err != nil {
		log.Fatalln(err)
	}

	// initiate image server
	imgServ, err := NewStaticImageServer("data")
	if err != nil {
		log.Fatalln(err)
	}

	// create errgroup to capture the error in parallel tasks
	g, _ := errgroup.WithContext(context.Background())

	// run telegram bot
	g.Go(func() error {
		return RunTelegramBot(cfg.TelegramToken, imgServ)
	})

	// run discord bot
	g.Go(func() error {
		return RunDiscordBot(cfg.DiscordToken, imgServ)
	})

	// wait for earliest error captured
	if err := g.Wait(); err != nil {
		log.Fatalln(err)
	}
}
