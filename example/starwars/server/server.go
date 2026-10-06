package main

import (
	"log"
	"net/http"

	"github.com/tribunadigital/graphql-go"
	"github.com/tribunadigital/graphql-go/example/internal/graphiql"
	"github.com/tribunadigital/graphql-go/example/starwars"
	"github.com/tribunadigital/graphql-go/relay"
)

func main() {
	schema := graphql.MustParseSchema(starwars.Schema, &starwars.Resolver{})

	http.Handle("GET /", graphiql.Handler())
	http.Handle("POST /query", &relay.Handler{Schema: schema})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
