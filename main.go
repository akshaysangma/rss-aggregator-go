package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/akshaysangma/rss-aggregator-go/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

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

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

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
	v1Router.Get("/readiness", handleReadiness)
	v1Router.Get("/err", handleErr)
	v1Router.Post("/user", apiCfg.handleUser)
	v1Router.Get("/user", apiCfg.middlewareAuth(apiCfg.handleGetUser))
	v1Router.Post("/feed", apiCfg.middlewareAuth(apiCfg.handleCreateFeed))

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
