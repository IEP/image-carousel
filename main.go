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

	cfg, err := LoadConfig("config")
	if err != nil {
		log.Fatalln(err)
	}

	imgServ, err := NewStaticImageServer("data")
	if err != nil {
		log.Fatalln(err)
	}

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return RunTelegramBot(cfg.TelegramToken, imgServ)
	})
	g.Go(func() error {
		return RunDiscordBot(cfg.DiscordToken, imgServ)
	})

	if err := g.Wait(); err != nil {
		log.Fatalln(err)
	}
}
