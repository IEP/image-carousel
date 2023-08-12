package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	tb "gopkg.in/telebot.v3"
)

var telegramTimeout = 1 * time.Minute

// RunTelegramBot with provided token and image server
func RunTelegramBot(token string, imgSrv ImageServer) error {
	// initate discord bot instance
	b, err := tb.NewBot(tb.Settings{
		Token: token,
		Poller: &tb.LongPoller{
			Timeout: 10 * time.Second,
		},
		Client: &http.Client{
			Timeout: telegramTimeout,
		},
	})
	if err != nil {
		return fmt.Errorf("tb.NewBot: %w", err)
	}
	err = b.RemoveWebhook()
	if err != nil {
		return fmt.Errorf("b.RemoveWebhook: %w", err)
	}

	buckets := imgSrv.GetBucketsName()
	commands := make([]tb.Command, 0)

	for _, bucket := range buckets {
		// add new command and handler
		bucket := bucket
		commands = append(commands, tb.Command{
			Text:        bucket,
			Description: fmt.Sprintf("Get random image from '%s' bucket", bucket),
		})
		b.Handle("/"+bucket, func(c tb.Context) error {
			// get random image based on the bucket name
			img := imgSrv.GetRandomImage(bucket)

			if img.Description == "" && img.PhotoPath == "" {
				log.Println("image not found")
				return errors.New("image not found")
			}

			log.Printf("[telegram] user: %+v", c.Sender())
			log.Printf("[telegram] image: %+v", img)

			if _, err := b.Reply(c.Message(), &tb.Photo{
				File:    tb.File{FileLocal: img.PhotoPath},
				Caption: img.Description,
			}); err != nil {
				log.Println("b.Reply:", err)
				return fmt.Errorf("b.Reply: %w", err)
			}

			return nil
		})
	}

	// set the command and start the bot
	if err := b.SetCommands(commands); err != nil {
		return fmt.Errorf("b.SetCommands: %w", err)
	}
	b.Start()

	return nil
}
