package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joaodddev/fraud-shield-go/internal/handler"
	"github.com/joaodddev/fraud-shield-go/internal/kafka"
	"github.com/joaodddev/fraud-shield-go/internal/rules"
)

func main() {
	broker := getEnv("KAFKA_BROKER", "localhost:9092")
	topic := getEnv("KAFKA_TOPIC_TRANSACTION_ANALYZED", "transaction.analyzed")
	port := getEnv("APP_PORT", "8080")

	producer, err := kafka.NewProducer(broker, topic)
	if err != nil {
		slog.Error("failed to connect kafka producer", "error", err)
		os.Exit(1)
	}
	defer producer.Close()

	engine := rules.NewEngine()
	txHandler := handler.NewTransactionHandler(engine, producer)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Get("/health", handler.Health)
	r.Post("/transactions", txHandler.Analyze)

	addr := fmt.Sprintf(":%s", port)
	slog.Info("server starting", "addr", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
