package main

import (
	"context"
	"encoding/json"
	otel "github.com/luizalexandresilva/goexpert-labs-otel/pkg/otel"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	chSo := make(chan os.Signal, 1)
	signal.Notify(chSo, os.Interrupt, syscall.SIGINT)

	ctx, shutdownSo := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer shutdownSo()

	ShutdownProvider, err := otel.InitProvider(ctx, "service-a", "otel-collector:4317")
	if err != nil {
		slog.Error("[InitProvider]", "error", err.Error())
		os.Exit(5)
	}
	defer func() {
		if err := ShutdownProvider(ctx); err != nil {
			slog.Error("[ShutdownProvider]", "error", err.Error())
			os.Exit(5)
		}
	}()

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var dto struct {
		Cep string `json:"cep"`
	}
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validateCep(w, dto.Cep)
}

func validateCep(w http.ResponseWriter, cep string) {
	if len(cep) != 8 || reflect.TypeOf(cep).Kind() != reflect.String {
		http.Error(w, "Invalid zipcode", http.StatusUnprocessableEntity)
		return
	}
}
