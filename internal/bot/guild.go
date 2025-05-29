package bot

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func IntervalRefreshAll(s *discordgo.Session) {
	for {
		log.Printf("Fetching all guilds")
		ProcessAllGuilds(s)
		time.Sleep(time.Hour * 24)
	}
}
