package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/akshaysangma/rss-aggregator-go/handler"
	"github.com/akshaysangma/rss-aggregator-go/internal/database"
	"github.com/akshaysangma/rss-aggregator-go/internal/rss"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalln("PORT env variable not set")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalln("DB_URL env variable not set")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	db := database.New(conn)
	apiCfg := handler.ApiConfig{
		DB: db,
	}

	go rss.StartScrapper(db, 10, 10*time.Second)

	router := chi.NewRouter()

	// TODO : Make it more restrictive
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/readiness", handler.HandleReadiness)
	v1Router.Get("/err", handler.HandleErr)

	v1Router.Post("/user", apiCfg.HandleUser)
	v1Router.Get("/user", apiCfg.MiddlewareAuth(apiCfg.HandleGetUser))
	v1Router.Get("/user/posts", apiCfg.MiddlewareAuth(apiCfg.HandleGetPostForUser))

	v1Router.Post("/feed", apiCfg.MiddlewareAuth(apiCfg.HandleCreateFeed))
	v1Router.Get("/feeds", apiCfg.HandleGetFeeds)

	v1Router.Post("/feed_follows", apiCfg.MiddlewareAuth(apiCfg.HandleCreateFeedFollows))
	v1Router.Get("/feed_follows", apiCfg.MiddlewareAuth(apiCfg.HandleGetFeedFollowsByUser))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.MiddlewareAuth(apiCfg.HandleDeleteFeedFollows))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Println("Server running on port :", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
