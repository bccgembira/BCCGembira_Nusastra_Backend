package dto

type QuestionResponse struct {
	ID     int64  `json:"id"`
	QuizID int64  `json:"quiz_id"`
	Title  string `json:"title"`
}
