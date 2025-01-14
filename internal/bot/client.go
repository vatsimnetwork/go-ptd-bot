package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"log"
	"os"
	"os/signal"
	"ptd-discord-bot/commands"
	"ptd-discord-bot/commands/roles"
	"ptd-discord-bot/functions"
	"ptd-discord-bot/internal/config"
)

func Session() (*discordgo.Session, error) {
	discord, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		return nil, err
	}
	return discord, nil
}

func Run() {
	log.Print("Starting discord-bot-v2")
	s, err := Session()
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
		panic(err.Error())
	}
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	AddMemberHandlers(s)

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		go IntervalRefreshAll(s)
		go config.IntervalReloadConfigs()
		err := s.UpdateWatchStatus(0, "The VATSIM Network")
		if err != nil {
			sentry.CaptureException(err)
		}
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	s.AddHandler(func(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
		go functions.ProcessMember(s, e.GuildID, e.Member)
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var GuildCommandHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
			"member-roles": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				roles.HandleMemberRoles(s, i)
			},
		}
		if h, ok := GuildCommandHandler[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	err = s.Open()
	if err != nil {
		println(err.Error())
		panic(err.Error())

	}
	log.Println("Adding commands...")
	registeredGuildCommands := make([]*discordgo.ApplicationCommand, len(commands.GuildCommands))
	for i, v := range commands.GuildCommands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredGuildCommands[i] = cmd
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}
