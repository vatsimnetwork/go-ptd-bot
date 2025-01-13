package functions

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/carlmjohnson/requests"
	"github.com/getsentry/sentry-go"
	"ptd-discord-bot/internal/config"
	"reflect"
)

type CidResponse struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
}

type V2Response struct {
	Rating      int `json:"rating"`
	PilotRating int `json:"pilotrating"`
}

func getCID(m *discordgo.User) (*CidResponse, error) {

	var CID *CidResponse
	err := requests.
		URL("https://api.vatsim.net").
		Pathf("/v2/members/discord/%s", m.ID).
		CheckStatus(200).
		ToJSON(&CID).
		Fetch(context.Background())

	if requests.HasStatusErr(err, 404) {
		return nil, err
	} else if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}
	return CID, nil
}

func getRatings(u *discordgo.User) (*V2Response, error) {
	var RatingsResponse *V2Response

	m, err := getCID(u)
	if m == nil || err != nil {
		return nil, err
	}

	err = requests.
		URL("https://api.vatsim.net").
		Pathf("/v2/members/%s", m.UserId).
		ToJSON(&RatingsResponse).
		CheckStatus(200).
		Fetch(context.Background())

	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}
	return RatingsResponse, nil

}

func ProcessMember(s *discordgo.Session, g string, m *discordgo.User) {
	mem, err := getRatings(m)
	if mem == nil || err != nil {
		sentry.CaptureException(err)
		return
	}

	ratingRoles := config.GetRatingsRoles()
	pilotRatingRoles := config.GetPilotRatingRoles()

	//var rolesEmbed []string
	newRoles := new([]string)
	CurrentRoles := make(map[string]string)
	UpdatedRoles := make(map[string]string)

	mr, _ := s.GuildMember(g, m.ID)

	for _, v := range mr.Roles {
		CurrentRoles[v] = v
		UpdatedRoles[v] = v
	}

	for _, v := range ratingRoles.Ratings {
		if CurrentRoles[v.DisocrdRoleId] == v.DisocrdRoleId {
			delete(UpdatedRoles, v.DisocrdRoleId)
		}
		if mem.Rating == v.CertValue {
			UpdatedRoles[v.DisocrdRoleId] = v.DisocrdRoleId
		}
	}
	for _, v := range pilotRatingRoles.Ratings {
		if CurrentRoles[v.DisocrdRoleId] == v.DisocrdRoleId {
			delete(UpdatedRoles, v.DisocrdRoleId)
		}
		if mem.PilotRating == v.CertValue {
			UpdatedRoles[v.DisocrdRoleId] = v.DisocrdRoleId
		}
	}

	for k := range UpdatedRoles {
		*newRoles = append(*newRoles, k)
	}

	eqcheck := reflect.DeepEqual(CurrentRoles, UpdatedRoles)

	fmt.Println(eqcheck)

	if !eqcheck {
		_, err = s.GuildMemberEdit(g, m.ID, &discordgo.GuildMemberParams{
			Roles: newRoles,
		})
	}

	if err != nil {
		sentry.CaptureException(err)
		return
	}
}
