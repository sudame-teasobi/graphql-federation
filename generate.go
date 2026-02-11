//go:generate gqlgen generate --config gqlgen.user.yml
//go:generate gqlgen generate --config gqlgen.task.yml
//go:generate rover supergraph compose --config supergraph.yaml -o supergraph.graphql
package main
