package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"gft/internal/user"
	"gft/internal/user/graph"
	"gft/internal/user/model"
	"gft/internal/user/resolver"
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

const defaultPort = "8080"

//go:embed users.json
var initUsersJSON []byte

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	var initUsers []*model.User
	err := json.Unmarshal(initUsersJSON, &initUsers)
	if err != nil {
		panic(fmt.Errorf("failed to load init user data: %w", err))
	}

	repo := user.NewRepository(initUsers)

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
