package models

type Submission struct {
	ID        int
	UserID    int
	ProblemID int
	Code      string
	Result    string
}
