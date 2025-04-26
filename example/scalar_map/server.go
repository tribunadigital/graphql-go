package main

import (
	"fmt"
	"log"
	"net/http"

	graphql "github.com/tribunadigital/graphql-go"
	"github.com/tribunadigital/graphql-go/relay"
)

type Map map[string]interface{}

func (Map) ImplementsGraphQLType(name string) bool {
	return name == "Map"
}

func (m *Map) UnmarshalGraphQL(input interface{}) error {
	val, ok := input.(map[string]interface{})
	if !ok {
		return fmt.Errorf("wrong type")
	}
	*m = val
	return nil
}

type Args struct {
	Name string
	Data Map
}

type mutation struct{}

func (_ *mutation) Hello(args Args) string {

	fmt.Println(args)

	return "Args accept!"
}

func main() {
	s := `
		scalar Map

		type Query {}

		type Mutation {
			hello(
				name: String!
				data: Map!
			): String!
		}
	`
	schema := graphql.MustParseSchema(s, &mutation{})
	http.Handle("/query", &relay.Handler{Schema: schema})

	log.Println("Listen in port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
