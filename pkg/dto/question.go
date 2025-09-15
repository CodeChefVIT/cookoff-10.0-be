package dto

import "github.com/google/uuid"

type Question struct {
	ID               uuid.UUID `json:"id"`
	Description      string    `json:"description"`
	Title            string    `json:"title"`
	Qtype            string    `json:"qType"`
	Isbountyactive   bool      `json:"isBountyActive"`
	InputFormat      []string  `json:"inputFormat"`
	Points           int32     `json:"points"`
	Round            int32     `json:"round"`
	Constraints      []string  `json:"constraints"`
	OutputFormat     []string  `json:"outputFormat"`
	SampleTestInput  []string  `json:"sampleTestInput"`
	SampleTestOutput []string  `json:"sampleTestOutput"`
	Explanation      []string  `json:"explanation"`
}
