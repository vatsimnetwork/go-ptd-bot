package main

import (
	"github.com/getsentry/sentry-go"
	"log"
	"ptd-discord-bot/internal/bot"
	"ptd-discord-bot/internal/config"
	"time"
)

func main() {
	bot.Run()
	if config.Env == "production" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:         config.SentryDSN,
			Debug:       false,
			Environment: config.Env,
		})
		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
		defer sentry.Flush(2 * time.Second)
	}
}
