package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Falcon50012/WB-L2/18/internal/calendar"
	"github.com/Falcon50012/WB-L2/18/internal/config"
	"github.com/Falcon50012/WB-L2/18/internal/handlers"
	"github.com/Falcon50012/WB-L2/18/internal/middleware"
)

func main() {
	cfg := config.Load()

	logger := log.New(os.Stdout, "[calendar] ", log.LstdFlags|log.LUTC)

	cal := calendar.New()
	h := handlers.New(cal)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	logged := middleware.Logging(logger)(mux)

	addr := fmt.Sprintf(":%s", cfg.Port)
	logger.Printf("starting server on %s", addr)

	if err := http.ListenAndServe(addr, logged); err != nil {
		logger.Fatalf("server error: %v", err)
	}
}
