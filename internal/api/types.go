package api

type LinkedAccountResponse struct {
	DiscordID string `json:"id"`
	UserID    string `json:"user_id"`
}

type MemberResponse struct {
	Rating         int `json:"rating"`
	PilotRating    int `json:"pilotrating"`
	MilitaryRating int `json:"militaryrating"`
}
