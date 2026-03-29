package commands

import (
	"flag"

	"github.com/vatsimnetwork/go-ptd-bot/commands/helpers"
	"github.com/vatsimnetwork/go-ptd-bot/commands/roles"
	"github.com/vatsimnetwork/go-ptd-bot/internal/config"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate) error

func genEnvGuild() string {
	if config.Env == "dev" {
		return "1037908270737784872"
	} else if config.Env == "prod" {
		return "901078003482783765"
	}
	return ""
}

var (
	AdminPermissions int64 = discordgo.PermissionAdministrator
	GuildID                = flag.String("guild", genEnvGuild(), "Test guild ID. If not passed - bot registers commands globally")
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
		{
			Name:        "fetch-nmoc",
			Description: "Fetches NMOC Data such as time to take exam and how many attempts",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "cid",
					Description: "The Members CID to check",
					Required:    true,
				},
			},
			DefaultMemberPermissions: &AdminPermissions,
		},
	}
	GuildCommandHandlers = map[string]CommandHandler{
		"member-roles": roles.HandleMemberRoles,
		"fetch-nmoc":   helpers.NMOC,
	}
)
