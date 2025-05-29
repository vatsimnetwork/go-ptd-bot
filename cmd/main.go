package main

import (
	"log"
	"time"

	"github.com/vatsimnetwork/go-ptd-bot/internal/bot"
	"github.com/vatsimnetwork/go-ptd-bot/internal/config"

	"github.com/getsentry/sentry-go"
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
