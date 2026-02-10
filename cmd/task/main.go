package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"gft/internal/task"
	"gft/internal/task/graph"
	"gft/internal/task/graph/resolver"
	"gft/internal/task/model"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8081"

//go:embed tasks.json
var initTasksJSON []byte

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	var initTasks []*model.Task
	err := json.Unmarshal(initTasksJSON, &initTasks)
	if err != nil {
		panic(fmt.Errorf("failed to load init tasks: %w", err))
	}

	repo := task.NewRepository(initTasks)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver.Resolver{Repo: repo}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
