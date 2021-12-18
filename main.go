package main

import (
	"log"
	"math/rand"
	"time"
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

	if err := RunTelegramBot(cfg.TelegramToken, imgServ); err != nil {
		log.Fatalln(err)
	}
}
