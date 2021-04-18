package main

import (
	"net/http"
	"os"

	"go.uber.org/zap"

	"github.com/pcherednichenko/spam_fighter_bot/internal/app"
)

func main() {
	productionLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() { _ = productionLogger.Sync() }()
	sugaredLogger := productionLogger.Sugar()

	sugaredLogger.Info("Starting!")

	// service health check
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	go func() {
		err := http.ListenAndServe(":80", nil)
		if err != nil {
			sugaredLogger.Fatal("error while creating health handler", err)
		}
	}()

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		sugaredLogger.Fatal("bot token should be not empty")
	}
	app.StartBot(sugaredLogger, botToken)
}
