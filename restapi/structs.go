package restapi

type MoodleUsers struct {
	Users []MoodleUser `json:"users"`
}

type MoodleUser struct {
	ID int `json:"id"`
}

type ExamAttempts struct {
	ExamAttempt []ExamAttempt `json:"attempts"`
}

type ExamAttempt struct {
	Attempt   int     `json:"attempt"`
	IsPreview int     `json:"preview"`
	TimeStart int64   `json:"timestart"`
	TimeEnd   int64   `json:"timefinish"`
	SumGrade  float64 `json:"sumgrades"`
}
