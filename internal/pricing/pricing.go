package pricing

import (
	"fmt"
)

type pricingDetail struct {
	inputTokensPerMillion  float64
	outputTokensPerMillion float64
}

var pricingMetadata = map[string]pricingDetail{
	"gpt-4o-mini": {
		inputTokensPerMillion:  0.150,
		outputTokensPerMillion: 0.600,
	},
	"gpt-4o": {
		inputTokensPerMillion:  2.50,
		outputTokensPerMillion: 10.00,
	},
}

// Sample budget metadata

var budgetMetadata = map[string]float64{
	"us-1": 100,
	"us-2": 200,
	"us-3": 300,
	"us-4": 100,
	"us-5": 500,
	"us-6": 600,
}

func CalculateSpend(model string, inputTokens, outputTokens float64) (float64, error) {

	val, ok := pricingMetadata[model]
	if !ok {
		return -1, fmt.Errorf("model %s does not exist", model)
	}
	spend := (val.inputTokensPerMillion * inputTokens) + (val.outputTokensPerMillion * (outputTokens))
	// spend = math.Round(spend*10000) / 10000
	return spend, nil
}

func GetBudget(useCaseId string) float64 {
	budget, ok := budgetMetadata[useCaseId]
	if !ok {
		return -1
	}

	return budget
}
