package domain

import "time"

type RiskLevel string
type Decision string

const (
	RiskLow    RiskLevel = "LOW"
	RiskMedium RiskLevel = "MEDIUM"
	RiskHigh   RiskLevel = "HIGH"
)

const (
	DecisionApproved Decision = "APPROVED"
	DecisionBlocked  Decision = "BLOCKED"
)

type AnalysisResult struct {
	TransactionID string    `json:"transaction_id"`
	Score         int       `json:"score"`
	RiskLevel     RiskLevel `json:"risk_level"`
	Decision      Decision  `json:"decision"`
	Reasons       []string  `json:"reasons"`
	AnalyzedAt    time.Time `json:"analyzed_at"`
}
