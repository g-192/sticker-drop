package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/g-192/stickerdrop/graph"
)

const defaultPort = "8080"

func main() {
	ctx := context.Background()

	dbUrl := "postgres://mule:secretpassword@localhost:5432/stickerdrop"
	dbPool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatalf("Error connecting to Postgres: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Postgres not answering to Ping: %v", err)
	}
	fmt.Println("Successfully connected to PostgresSQL!")

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis not answering to Ping: %v", err)
	}
	fmt.Println("Successfully connected to Redis!")

	_, err = dbPool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS drops(
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			total_available INT NOT NULL,
			claimed INT NOT NULL DEFAULT 0
		);
		INSERT INTO drops (id, title, total_available, claimed)
		VALUES (1, '500 Stickers Giveaway', 500, 0)
		ON CONFLICT (id) DO NOTHING;
	`)
	if err != nil {
		log.Fatalf("Error creating the table: %v", err)
	}
	fmt.Println("Postgres table 'drops' is ready!")

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			DB:    dbPool,
			Redis: rdb,
		},
	}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	http.Handle("/query", corsMiddleware(srv))

	fmt.Printf("Server is running! Open your Browser under: http://localhost:%s/\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// CORS Middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // for development
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

		//GraphQL Clients send OPTIONS requests first to check for permission
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
