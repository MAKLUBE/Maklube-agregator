package main

import (
	"log"
	"net/http"
	"os"
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

	err = srv.ListenAndServe()
	logger.Fatal(err)
}

func getenv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
