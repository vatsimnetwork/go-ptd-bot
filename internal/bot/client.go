package bot

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/vatsimnetwork/go-ptd-bot/commands"
	"github.com/vatsimnetwork/go-ptd-bot/internal/config"
	"github.com/vatsimnetwork/go-ptd-bot/internal/util"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
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
	s, err := Session()
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalf("failed to create discord session: %v", err)
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
		go util.ProcessMember(s, e.GuildID, e.Member)
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commands.GuildCommandHandlers[i.ApplicationCommandData().Name]; ok {
			if err := h(s, i); err != nil {
				sentry.CaptureException(err)
			}
		}
	})

	if err := s.Open(); err != nil {
		log.Fatalf("failed to connect to discord: %v", err)
	}

	defer func(s *discordgo.Session) {
		if err := s.Close(); err != nil {
			log.Fatalf("failed to close discord session: %v", err)
		}
	}(s)

	log.Println("Adding commands...")
	for _, v := range commands.GuildCommands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}
