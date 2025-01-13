package commands

import (
	"github.com/bwmarrin/discordgo"
)

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
)
