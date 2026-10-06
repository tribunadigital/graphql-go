package main

import (
	"log"
	"net/http"

	"github.com/tribunadigital/graphql-go"
	"github.com/tribunadigital/graphql-go/example/internal/graphiql"
	"github.com/tribunadigital/graphql-go/example/social"
	"github.com/tribunadigital/graphql-go/relay"
)

func main() {
	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers(), graphql.MaxParallelism(20)}
	schema := graphql.MustParseSchema(social.Schema, &social.Resolver{}, opts...)

	http.Handle("GET /", graphiql.Handler())
	http.Handle("POST /query", &relay.Handler{Schema: schema})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
