package api

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/carlmjohnson/requests"
)

func GetLinkedAccount(m *discordgo.User) (*LinkedAccountResponse, error) {
	var res *LinkedAccountResponse

	err := requests.
		URL("https://api.vatsim.net").
		Pathf("/v2/members/discord/%s", m.ID).
		CheckStatus(200).
		ToJSON(&res).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetMember(u *discordgo.User) (*MemberResponse, error) {
	var res *MemberResponse

	m, err := GetLinkedAccount(u)
	if m == nil || err != nil {
		return nil, err
	}

	err = requests.
		URL("https://api.vatsim.net").
		Pathf("/v2/members/%s", m.UserID).
		ToJSON(&res).
		CheckStatus(200).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}

	return res, nil
}
