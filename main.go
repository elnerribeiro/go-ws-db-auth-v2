package main

import (
	"context"
	"errors"
	"github.com/elnerribeiro/go-ws-db-auth-v2/db"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/elnerribeiro/go-ws-db-auth-v2/app"
	"github.com/elnerribeiro/go-ws-db-auth-v2/controllers"
	"github.com/elnerribeiro/go-ws-db-auth-v2/utils"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	logger, _ := utils.GetLoggerAndContext()

	router := mux.NewRouter()

	router.HandleFunc("/api/users", controllers.ListUsers).Methods("POST")
	router.HandleFunc("/api/user/{id:[0-9]+}", controllers.GetUserByID).Methods("GET")
	router.HandleFunc("/api/user", controllers.Upsert).Methods("PUT")
	router.HandleFunc("/api/user/{id:[0-9]+}", controllers.Delete).Methods("DELETE")
	router.HandleFunc("/api/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/validate", controllers.Validate).Methods("GET")
	router.HandleFunc("/api/insert/{id:[0-9]+}", controllers.ListInsert).Methods("GET")
	router.HandleFunc("/api/insert/sync/{qty:[0-9]+}", controllers.InsertSync).Methods("PUT")
	router.HandleFunc("/api/insert/async/{qty:[0-9]+}", controllers.InsertASync).Methods("PUT")
	router.HandleFunc("/api/insert", controllers.ClearInserts).Methods("DELETE")
	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Error().Msgf("[NotFoundHandler] Resource not found: %s", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		Debug:            true,
	})

	handler := c.Handler(router)

	srv := &http.Server{
		Addr:    ":8000",
		Handler: handler,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error().Err(err).Msgf("[main] Error while executing server: %s", err)
		}
	}()

	logger.Info().Msg("[main] Server started on port 8000")

	<-done
	logger.Info().Msg("[main] Server Stopped")
	db.FinalizeDB(logger)

	endContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(endContext); err != nil {
		logger.Error().Err(err).Msgf("[main] Server Shutdown Failed:%+v", err)
	}
	logger.Info().Msg("[main] Server Exited Properly")
}
