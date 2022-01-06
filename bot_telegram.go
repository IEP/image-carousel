package main

import (
	"fmt"
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

// RunTelegramBot with provided token and image server
func RunTelegramBot(token string, imgSrv ImageServer) error {
	// initate discord bot instance
	b, err := tb.NewBot(tb.Settings{
		Token: token,
		Poller: &tb.LongPoller{
			Timeout: 10 * time.Second,
		},
	})
	if err != nil {
		return fmt.Errorf("tb.NewBot: %w", err)
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
		b.Handle("/"+bucket, func(m *tb.Message) {
			// get random image based on the bucket name
			img := imgSrv.GetRandomImage(bucket)

			if img.Description == "" && img.PhotoPath == "" {
				log.Println("image not found")
				return
			}

			log.Printf("[telegram] user: %+v", m.Sender)
			log.Printf("[telegram] image: %+v", img)

			if _, err := b.Reply(m, &tb.Photo{
				File:    tb.File{FileLocal: img.PhotoPath},
				Caption: img.Description,
			}); err != nil {
				log.Println("b.Reply:", err)
			}
		})
	}

	// set the command and start the bot
	_ = b.SetCommands(commands)
	b.Start()

	return nil
}
