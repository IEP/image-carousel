package main

import (
	"context"
	"log"
	"os"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
)

// RunDiscordBot instance with provided token and image server
func RunDiscordBot(token string, imgSrv ImageServer) error {
	// get buckets information
	buckets := imgSrv.GetBucketsName()
	bucketFilter := make(map[string]bool)
	for _, bucketName := range buckets {
		bucketFilter[bucketName] = true
	}

	const prefix = "!"

	// initiate discord bot instance
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

	// add filter
	filter, _ := std.NewMsgFilter(context.Background(), client)
	filter.SetPrefix(prefix)

	// add message handler
	client.Gateway().WithMiddleware(
		filter.NotByBot,
		filter.HasPrefix,
		filter.StripPrefix,
	).MessageCreate(func(s disgord.Session, data *disgord.MessageCreate) {
		msg := data.Message

		// check whether the requested bucket is exists
		if bucketFilter[msg.Content] {
			img := imgSrv.GetRandomImage(msg.Content)

			if img.Description == "" && img.PhotoPath == "" {
				log.Println("image not found")
				return
			}

			log.Printf("[discord] user: %+v", data.Message.Author)
			log.Printf("[discord] image: %+v", img)

			// get file that will be fed as io.Reader in Files
			f, err := os.Open(img.PhotoPath)
			if err != nil {
				log.Println("os.Open:", err)
			}
			defer f.Close()

			// send message with embedded image
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
