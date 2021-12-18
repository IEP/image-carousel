package main

import (
	"fmt"
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func RunTelegramBot(token string, imgSrv ImageServer) error {
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
		bucket := bucket
		commands = append(commands, tb.Command{
			Text:        bucket,
			Description: fmt.Sprintf("Get random image from '%s' bucket", bucket),
		})
		b.Handle("/"+bucket, func(m *tb.Message) {
			img := imgSrv.GetRandomImage(bucket)

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

	_ = b.SetCommands(commands)
	b.Start()

	return nil
}
