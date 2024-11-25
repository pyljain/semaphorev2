package entities

type Total struct {
	UseCaseId string  `json:"useCaseId" validate:"required" bson:"_id"`
	Spend     float64 `json:"spend" bson:"spend"`
}
