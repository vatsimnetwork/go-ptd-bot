package helpers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/vatsimnetwork/go-ptd-bot/internal/config"
	"github.com/vatsimnetwork/go-ptd-bot/restapi"
)

func NMOC(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	options := i.ApplicationCommandData().Options

	user := options[0].IntValue()

	var uo restapi.MoodleUsers
	var at restapi.ExamAttempts
	var fo string = ""
	o, e := restapi.SendRequest("GET", fmt.Sprintf("https://learn.vatsim.net/webservice/rest/server.php?wstoken=%v&moodlewsrestformat=json&wsfunction=core_user_get_users&criteria[0][key]=username&criteria[0][value]=%v", config.MoodleToken, user), nil)
	if e != nil {
		return e
	}
	err := json.NewDecoder(o.Body).Decode(&uo)
	if err != nil {
		return err
	}
	oa, ea := restapi.SendRequest("GET", fmt.Sprintf("https://learn.vatsim.net/webservice/rest/server.php?wstoken=%v&moodlewsrestformat=json&wsfunction=mod_quiz_get_user_attempts&quizid=30&userid=%v&status=all", config.MoodleToken, uo.Users[0].ID), nil)
	if ea != nil {
		fmt.Println(e)
	}
	err = json.NewDecoder(oa.Body).Decode(&at)
	if len(at.ExamAttempt) == 0 {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No exam attempts found for user",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return nil
	}
	if err != nil {
		return err
	}

	for i := 0; i < len(at.ExamAttempt); i++ {
		fo += fmt.Sprintf("Attempt: %v\n", at.ExamAttempt[i].Attempt)
		fo += fmt.Sprintf("Time Start: %v\n", time.Unix(at.ExamAttempt[i].TimeStart, 0))
		fo += fmt.Sprintf("Time End: %v\n", time.Unix(at.ExamAttempt[i].TimeEnd, 0))
		st := time.Unix(at.ExamAttempt[0].TimeEnd, 0)
		tott := st.Sub(time.Unix(at.ExamAttempt[0].TimeStart, 0))
		fo += fmt.Sprintf("Exam took %v\n", tott)
		math := at.ExamAttempt[i].SumGrade / float64(40) * 100
		fo += fmt.Sprintf("Grade: %v\n", math)
		tpq := tott / 40
		fo += fmt.Sprintf("Time Per Question: %v\n", tpq)
		fo += "-------------------------------------------------\n"
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fo,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
