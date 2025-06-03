package util

import (
	"fmt"
	"reflect"

	"github.com/vatsimnetwork/go-ptd-bot/internal/api"
	"github.com/vatsimnetwork/go-ptd-bot/internal/config"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
)

func ProcessMember(s *discordgo.Session, guildID string, m *discordgo.Member) {
	mem, err := api.GetMember(m.User)
	if mem == nil {
		return
	} else if err != nil {
		sentry.CaptureException(err)
		return
	}

	ratingRoles := config.GetRatingsRoles(guildID)
	pilotRatingRoles := config.GetPilotRatingRoles(guildID)
	militaryRatingRoles := config.GetMilitaryRatingRoles(guildID)

	actualRoles := make(map[string]struct{}, len(m.Roles))
	expectedRoles := make(map[string]struct{}, len(m.Roles))

	for _, v := range m.Roles {
		actualRoles[v] = struct{}{}
		expectedRoles[v] = struct{}{}
	}

	for _, v := range ratingRoles {
		if _, ok := expectedRoles[v.RoleID]; ok {
			delete(expectedRoles, v.RoleID)
		}
		if mem.Rating == v.CertValue {
			expectedRoles[v.RoleID] = struct{}{}
		}
	}
	for _, v := range pilotRatingRoles {
		if _, ok := expectedRoles[v.RoleID]; ok {
			delete(expectedRoles, v.RoleID)
		}
		if mem.PilotRating == v.CertValue {
			expectedRoles[v.RoleID] = struct{}{}
		}
	}
	for _, v := range militaryRatingRoles {
		if _, ok := expectedRoles[v.RoleID]; ok {
			delete(expectedRoles, v.RoleID)
		}
		if mem.MilitaryRating == v.CertValue {
			expectedRoles[v.RoleID] = struct{}{}
		}
	}

	if reflect.DeepEqual(actualRoles, expectedRoles) {
		return
	}

	fmt.Println(actualRoles, expectedRoles)

	roleIDs := new([]string)
	for k := range expectedRoles {
		*roleIDs = append(*roleIDs, k)
	}

	_, err = s.GuildMemberEdit(guildID, m.User.ID, &discordgo.GuildMemberParams{
		Roles: roleIDs,
	})
	if err != nil {
		sentry.CaptureException(err)
	}
}
