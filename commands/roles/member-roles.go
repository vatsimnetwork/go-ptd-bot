package roles

import (
	"github.com/vatsimnetwork/go-ptd-bot/internal/util"

	"github.com/bwmarrin/discordgo"
)

func HandleMemberRoles(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := i.ApplicationCommandData().Options

	user := options[0].UserValue(s)

	member, _ := s.GuildMember(i.GuildID, user.ID)

	go util.ProcessMember(s, i.GuildID, member)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Roles Assigned!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
