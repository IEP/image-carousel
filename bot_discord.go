package main

import (
	"context"
	"log"
	"os"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
)

func RunDiscordBot(token string, imgSrv ImageServer) error {
	buckets := imgSrv.GetBucketsName()
	bucketFilter := make(map[string]bool)
	for _, bucketName := range buckets {
		bucketFilter[bucketName] = true
	}

	const prefix = "!"

	client := disgord.New(disgord.Config{
		BotToken: token,
		RejectEvents: []string{
			disgord.EvtTypingStart,
			disgord.EvtPresenceUpdate,
			disgord.EvtGuildMemberAdd,
			disgord.EvtGuildMemberUpdate,
			disgord.EvtGuildMemberRemove,
		},
		DMIntents: disgord.IntentDirectMessages |
			disgord.IntentDirectMessageReactions |
			disgord.IntentDirectMessageTyping,
	})

	defer func() {
		_ = client.Gateway().Connect()
	}()

	filter, _ := std.NewMsgFilter(context.Background(), client)
	filter.SetPrefix(prefix)

	client.Gateway().WithMiddleware(
		filter.NotByBot,
		filter.HasPrefix,
		filter.StripPrefix,
	).MessageCreate(func(s disgord.Session, data *disgord.MessageCreate) {
		msg := data.Message

		if bucketFilter[msg.Content] {
			img := imgSrv.GetRandomImage(msg.Content)

			log.Printf("[discord] user: %+v", data.Message.Author)
			log.Printf("[discord] image: %+v", img)

			f, err := os.Open(img.PhotoPath)
			if err != nil {
				log.Println("os.Open:", err)
			}
			defer f.Close()

			if _, err := msg.Reply(context.Background(), s, &disgord.CreateMessageParams{
				Content: img.Description,
				Files: []disgord.CreateMessageFileParams{
					{Reader: f, FileName: "nft.jpg", SpoilerTag: false},
				},
			}); err != nil {
				log.Println("msg.Reply:", err)
			}
		}
	})

	return nil
}
