package main

import (
	"context"
	"log"

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
			if _, err := msg.Reply(context.Background(), s, msg.Content); err != nil {
				log.Println("msg.Reply", err)
			}
		}
	})

	client.Gateway().BotReady(func() {
		log.Println("Discord bot is ready!")
	})

	return nil
}
