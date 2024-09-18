package functions

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/carlmjohnson/requests"
	"log"
)

type discordUser struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
}

type vatsimUser struct {
	Id             int `json:"id"`
	Rating         int `json:"rating"`
	PilotRating    int `json:"pilotrating"`
	MilitaryRating int `json:"militaryrating"`
}

func ProcessMember(m *discordgo.Member) {

	ManagedRoles := [14]string{"Supervisor", "Administrator", "New Member", "PPL", "IR", "CMEL", "ATPL", "Flight Instructor", "Flight Examiner", "No Military Rating", "M1", "M2", "M3", "M4"}

	var user discordUser
	url := fmt.Sprintf("https://api.vatsim.net/v2/members/discord/%v", m.User.ID)
	err := requests.
		URL(url).
		Header("User-Agent", "PTDDiscordBotv2").
		ToJSON(&user).
		Fetch(context.Background())

	if err != nil {
		log.Println(err)
		return
	}

	var VatsimUser vatsimUser
	memberUrl := fmt.Sprintf("https://api.vatsim.net/v2/members/%v", user.UserId)
	errr := requests.
		URL(memberUrl).
		Header("User-Agent", "PTDDiscordBotv2").
		ToJSON(&VatsimUser).
		Fetch(context.Background())

	if errr != nil {
		log.Println(errr)
		return
	}

	var roles []string

	switch {

	}
}
