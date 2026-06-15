package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/joaodddev/fraud-shield-go/internal/domain"
	"github.com/joaodddev/fraud-shield-go/internal/kafka"
	"github.com/joaodddev/fraud-shield-go/internal/rules"
)

type TransactionHandler struct {
	engine   *rules.Engine
	producer *kafka.Producer
}

func NewTransactionHandler(engine *rules.Engine, producer *kafka.Producer) *TransactionHandler {
	return &TransactionHandler{
		engine:   engine,
		producer: producer,
	}
}

type createTransactionRequest struct {
	AccountID       string  `json:"account_id"`
	Amount          float64 `json:"amount"`
	MerchantID      string  `json:"merchant_id"`
	MerchantCountry string  `json:"merchant_country"`
}

type analyzeResponse struct {
	Transaction domain.Transaction    `json:"transaction"`
	Result      domain.AnalysisResult `json:"result"`
}

func (h *TransactionHandler) Analyze(w http.ResponseWriter, r *http.Request) {
	var req createTransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.AccountID == "" || req.Amount <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "account_id and amount are required")
		return
	}

	t := domain.Transaction{
		ID:              uuid.NewString(),
		AccountID:       req.AccountID,
		Amount:          req.Amount,
		MerchantID:      req.MerchantID,
		MerchantCountry: req.MerchantCountry,
		CreatedAt:       time.Now().UTC(),
	}

	result := h.engine.Analyze(t)

	event := kafka.NewTransactionAnalyzedEvent(t, result)
	if err := h.producer.Publish(t.ID, event); err != nil {
		slog.Error("failed to publish event", "error", err, "transaction_id", t.ID)
		writeError(w, http.StatusInternalServerError, "failed to process transaction")
		return
	}

	status := http.StatusOK
	if result.Decision == domain.DecisionBlocked {
		status = http.StatusForbidden
	}

	writeJSON(w, status, analyzeResponse{
		Transaction: t,
		Result:      result,
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
