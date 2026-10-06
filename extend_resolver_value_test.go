package graphql_test

import (
	"testing"

	"github.com/tribunadigital/graphql-go"
	"github.com/tribunadigital/graphql-go/gqltesting"
)

// extendCountry has the memory layout of extendBaseCountry and adds a field,
// the way api-gateway extends statistics resolvers (statCountry and others).
type extendBaseCountry struct {
	Name string
}

type extendCountry struct {
	extendBaseCountry
}

func (c *extendCountry) Iso() string { return "iso:" + c.Name }

// extendBroadcasters exposes the extended type as a struct value field and as a
// method returning a struct value: both reach the executor as non-addressable values.
type extendBroadcasters struct {
	Country extendBaseCountry
}

func (b extendBroadcasters) Home() extendBaseCountry { return b.Country }

type extendValueQuery struct{}

func (extendValueQuery) Broadcasters() []extendBroadcasters {
	return []extendBroadcasters{
		{Country: extendBaseCountry{Name: "ukraine"}},
		{Country: extendBaseCountry{Name: "spain"}},
	}
}

func TestExtendResolverOnValueFields(t *testing.T) {
	t.Parallel()

	gqltesting.RunTests(t, []*gqltesting.Test{
		{
			Schema: graphql.MustParseSchema(`
				schema {
					query: Query
				}

				type Query {
					broadcasters: [Broadcasters!]!
				}

				type Broadcasters {
					country: Country!
					home: Country!
				}

				type Country {
					name: String!
					iso: String!
				}`,
				&extendValueQuery{},
				graphql.UseFieldResolvers(),
				graphql.UseExtendResolver(map[string]any{"Country": &extendCountry{}}),
			),
			Query: `
				{
					broadcasters {
						country { name iso }
						home { name iso }
					}
				}
			`,
			ExpectedResult: `
				{
					"broadcasters": [
						{"country": {"name": "ukraine", "iso": "iso:ukraine"}, "home": {"name": "ukraine", "iso": "iso:ukraine"}},
						{"country": {"name": "spain", "iso": "iso:spain"}, "home": {"name": "spain", "iso": "iso:spain"}}
					]
				}
			`,
		},
	})
}
