package bot

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
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
	session, err := Session()
	if err != nil {
		println(err.Error())
	}
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	AddMemberHandlers(session)

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		go IntervalRefreshAll(s)
		s.UpdateWatchStatus(0, "The VATSIM Network")
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	session.AddHandler(func(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
		go functions.ProcessMember(s, e.GuildID, e.Member)
	})

	err = session.Open()
	if err != nil {
		println(err.Error())
		panic(err.Error())

	}
	defer session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}
