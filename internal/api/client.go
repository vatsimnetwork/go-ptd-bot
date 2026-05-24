package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"github.com/vatsimnetwork/go-ptd-bot/internal/config"
)

func GetLinkedAccount(m *discordgo.User) (la *LinkedAccountResponse, e error) {
	var bodyReader io.Reader

	req, err := http.NewRequest("GET", config.APIURL+"/v2/members/discord/"+m.ID, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "PTDDiscordBot/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&la)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return la, nil
}

func GetMember(u *discordgo.User) (mr *MemberResponse, e error) {

	m, err := GetLinkedAccount(u)
	if m == nil || err != nil {
		sentry.CaptureException(err)
		return nil, err
	}

	var bodyReader io.Reader

	req, err := http.NewRequest("GET", config.APIURL+"/v2/members/"+m.UserID, bodyReader)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "PTDDiscordBot/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&mr)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return mr, nil
}
