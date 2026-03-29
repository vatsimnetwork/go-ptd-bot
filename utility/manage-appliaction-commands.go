package utility

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"github.com/vatsimnetwork/go-ptd-bot/commands"
)

func HandleApplicationCommandUpdates(session *discordgo.Session) {
	registeredGuildCommands := make([]*discordgo.ApplicationCommand, len(commands.GuildCommands))
	log.Println("Updating Commands")
	for i, v := range commands.GuildCommands {
		registeredGuildCommands[i] = v
	}
	_, err := session.ApplicationCommandBulkOverwrite(session.State.User.ID, *commands.GuildID, registeredGuildCommands)
	if err != nil {
		sentry.CaptureException(err)
		println(err.Error())
	}
}
