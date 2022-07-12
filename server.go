package main

import (
	"github.com/garixx/howtographql/internal/auth"
	_ "github.com/garixx/howtographql/internal/auth"
	database "github.com/garixx/howtographql/internal/pkg/db/mysql"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/garixx/howtographql/graph"
	"github.com/garixx/howtographql/graph/generated"
	"github.com/go-chi/chi"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	router.Use(auth.Middleware)
	//router.Use(auth.Middleware())

	database.InitDb()
	defer database.CloseDB()
	database.Migrate()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
