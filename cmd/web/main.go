package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MAKLUBE/AP1_Final_Project/internal/app"
	"github.com/MAKLUBE/AP1_Final_Project/internal/db"
	"github.com/MAKLUBE/AP1_Final_Project/internal/handlers"
	"github.com/MAKLUBE/AP1_Final_Project/internal/middleware"
	"github.com/MAKLUBE/AP1_Final_Project/internal/storage/mongo"
)

func main() {
	logger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	mongoURI := getenv("MONGO_URI", "mongodb://localhost:27017")
	mongoDBName := getenv("MONGO_DB", "Makluber")

	// Mongo connect
	client, err := db.NewMongoClient(mongoURI, 10*time.Second)
	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		_ = client.Disconnect(nil)
	}()
	database := client.Database(mongoDBName)

	tc, err := app.NewTemplateCache("./ui/html")
	if err != nil {
		logger.Fatal(err)
	}

	userStore := mongo.NewUserStore(database)
	sessionStore := mongo.NewSessionStore(database)
	restaurantStore := mongo.NewRestaurantStore(database)
	menuStore := mongo.NewMenuItemStore(database)
	reviewStore := mongo.NewReviewStore(database)
	orderStore := mongo.NewOrderStore(database)
	statusStore := mongo.NewOrderStatusHistoryStore(database)
	halalStore := mongo.NewHalalVerificationStore(database)

	application := &app.Application{
		Logger:      logger,
		Templates:   tc,
		Users:       userStore,
		Sessions:    sessionStore,
		Restaurants: restaurantStore,
		MenuItems:   menuStore,
		Reviews:     reviewStore,
		Orders:      orderStore,
		OrderStatus: statusStore,
		Halal:       halalStore,
	}

	r := handlers.Routes(application)
	chain := middleware.New(application).Chain(r)

	addr := getenv("ADDR", ":4000")
	logger.Printf("starting server on %s", addr)

	srv := &http.Server{
		Addr:              addr,
		Handler:           chain,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	done := make(chan struct{})

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	go func() {
		defer close(done)
		cleanupCount := 0
		logger.Println("background: session cleanup worker started")

		for {
			select {
			case <-ticker.C:
				cleanupCount++
				logger.Printf("background: cleanup #%d started", cleanupCount)
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				err := application.Sessions.DeleteExpired(ctx, time.Now().UTC())
				cancel()
				if err != nil {
					logger.Printf("background: cleanup #%d failed: %v", cleanupCount, err)
				} else {
					logger.Printf("background: cleanup #%d completed successfully", cleanupCount)
				}

			case <-shutdown:
				logger.Println("background: received shutdown signal, stop worker")
				return
			}
		}
	}()

	serverErrors := make(chan error, 1)
	go func() {
		logger.Println("server: HTTP server starting")
		serverErrors <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		logger.Fatalf("server: error starting server: %v", err)

	case sig := <-shutdown:
		logger.Printf("server: received shutdown signal: %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		logger.Println("server: shutting down HTTP server")
		if err := srv.Shutdown(ctx); err != nil {
			logger.Printf("server: error during shutdown: %v", err)
			_ = srv.Close()
		}

		logger.Println("server: waiting for background worker to finish")
		<-done
		logger.Println("server: background worker stopped")

		logger.Println("server: shutdown complete")
	}
}

func getenv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
