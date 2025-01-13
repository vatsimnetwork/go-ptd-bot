package bot

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

func IntervalRefreshAll(s *discordgo.Session) {
	for {
		log.Printf("Fetching all guilds")
		ProcessAllGuilds(s)
		time.Sleep(time.Hour * 24)
	}
}
