package roles

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"ptd-discord-bot/functions"
)

func HandleMemberRoles(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := i.ApplicationCommandData().Options

	user := options[0].UserValue(s)

	member, _ := s.GuildMember(i.GuildID, user.ID)

	go functions.ProcessMember(s, i.GuildID, member)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Roles Assigned!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		sentry.CaptureException(err)
	}
}
