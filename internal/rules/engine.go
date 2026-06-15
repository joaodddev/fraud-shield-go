package rules

import (
	"time"

	"github.com/joaodddev/fraud-shield-go/internal/domain"
)

var suspiciousCountries = map[string]bool{
	"KP": true, // North Korea
	"IR": true, // Iran
	"CU": true, // Cuba
	"SY": true, // Syria
}

type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Analyze(t domain.Transaction) domain.AnalysisResult {
	score := 0
	reasons := []string{}

	// Regra 1: valor alto
	if t.Amount > 10000 {
		score += 40
		reasons = append(reasons, "transaction amount exceeds $10,000")
	} else if t.Amount > 5000 {
		score += 20
		reasons = append(reasons, "transaction amount exceeds $5,000")
	}

	// Regra 2: país suspeito
	if suspiciousCountries[t.MerchantCountry] {
		score += 40
		reasons = append(reasons, "merchant located in high-risk country")
	}

	// Regra 3: horário suspeito (00h–05h UTC)
	hour := t.CreatedAt.UTC().Hour()
	if hour >= 0 && hour < 5 {
		score += 20
		reasons = append(reasons, "transaction initiated during high-risk hours (00h-05h UTC)")
	}

	// Regra 4: merchant ID vazio
	if t.MerchantID == "" {
		score += 30
		reasons = append(reasons, "missing merchant identifier")
	}

	// Cap no score
	if score > 100 {
		score = 100
	}

	riskLevel := resolveRiskLevel(score)
	decision := resolveDecision(score)

	return domain.AnalysisResult{
		TransactionID: t.ID,
		Score:         score,
		RiskLevel:     riskLevel,
		Decision:      decision,
		Reasons:       reasons,
		AnalyzedAt:    time.Now().UTC(),
	}
}

func resolveRiskLevel(score int) domain.RiskLevel {
	switch {
	case score >= 70:
		return domain.RiskHigh
	case score >= 40:
		return domain.RiskMedium
	default:
		return domain.RiskLow
	}
}

func resolveDecision(score int) domain.Decision {
	if score >= 70 {
		return domain.DecisionBlocked
	}
	return domain.DecisionApproved
}
