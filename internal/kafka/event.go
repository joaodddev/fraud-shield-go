package kafka

import "github.com/joaodddev/fraud-shield-go/internal/domain"

type TransactionAnalyzedEvent struct {
	Transaction domain.Transaction    `json:"transaction"`
	Result      domain.AnalysisResult `json:"result"`
}

func NewTransactionAnalyzedEvent(t domain.Transaction, r domain.AnalysisResult) TransactionAnalyzedEvent {
	return TransactionAnalyzedEvent{
		Transaction: t,
		Result:      r,
	}
}
