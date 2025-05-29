package commands

import (
	"github.com/vatsimnetwork/go-ptd-bot/commands/roles"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate) error

var (
	AdminPermissions int64 = discordgo.PermissionAdministrator
	GuildCommands          = []*discordgo.ApplicationCommand{
		{
			Name:        "member-roles",
			Description: "Assigns a members roles",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "member",
					Description: "Member role assignment",
					Required:    true,
				},
			},
			DefaultMemberPermissions: &AdminPermissions,
		},
	}
	GuildCommandHandlers = map[string]CommandHandler{
		"member-roles": roles.HandleMemberRoles,
	}
)
