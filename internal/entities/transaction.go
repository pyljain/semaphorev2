package entities

import "time"

type Transaction struct {
	RequestId    string    `json:"requestId" validate:"required" bson:"_id"`
	Model        string    `json:"model" validate:"required" bson:"model"`
	InputTokens  int       `json:"inputTokens" validate:"required" bson:"inputTokens"`
	OutputTokens int       `json:"outputTokens" validate:"required" bson:"outputTokens"`
	UseCaseId    string    `json:"useCaseId" validate:"required" bson:"useCaseId"`
	Spend        float64   `json:"spend" bson:"spend"`
	Date         time.Time `json:"date" bson:"date"`
}
