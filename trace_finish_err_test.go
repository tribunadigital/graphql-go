package graphql_test

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/tribunadigital/graphql-go"
	gqlerrors "github.com/tribunadigital/graphql-go/errors"
	"github.com/tribunadigital/graphql-go/introspection"
)

type fieldErrRecordingTracer struct {
	mu       sync.Mutex
	fieldErr map[string]*gqlerrors.QueryError
}

func (t *fieldErrRecordingTracer) TraceQuery(ctx context.Context, queryString string, operationName string, variables map[string]any, varTypes map[string]*introspection.Type) (context.Context, func([]*gqlerrors.QueryError)) {
	return ctx, func([]*gqlerrors.QueryError) {}
}

func (t *fieldErrRecordingTracer) TraceField(ctx context.Context, label, typeName, fieldName string, trivial bool, args map[string]any) (context.Context, func(*gqlerrors.QueryError)) {
	return ctx, func(err *gqlerrors.QueryError) {
		t.mu.Lock()
		defer t.mu.Unlock()
		if t.fieldErr == nil {
			t.fieldErr = make(map[string]*gqlerrors.QueryError)
		}
		t.fieldErr[typeName+"."+fieldName] = err
	}
}

type traceErrBaseCountry struct {
	Name string
}

type traceErrCountry struct {
	traceErrBaseCountry
}

func (c *traceErrCountry) Iso() (*string, error) {
	return nil, errors.New("iso failed for " + c.Name)
}

type traceErrExtendedQuery struct{}

func (traceErrExtendedQuery) Country() *traceErrBaseCountry {
	return &traceErrBaseCountry{Name: "ukraine"}
}

type traceErrPlainQuery struct{}

func (traceErrPlainQuery) Hello() (*string, error) {
	return nil, errors.New("hello failed")
}

func TestTraceFieldFinishReceivesResolverError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		schema   string
		resolver any
		opts     []graphql.SchemaOpt
		query    string
		field    string
		wantPath []any
	}{
		{
			name: "extended resolver",
			schema: `
				schema {
					query: Query
				}

				type Query {
					country: Country!
				}

				type Country {
					name: String!
					iso: String
				}`,
			resolver: &traceErrExtendedQuery{},
			opts: []graphql.SchemaOpt{
				graphql.UseFieldResolvers(),
				graphql.UseExtendResolver(map[string]any{"Country": &traceErrCountry{}}),
			},
			query:    `{ country { name iso } }`,
			field:    "Country.iso",
			wantPath: []any{"country", "iso"},
		},
		{
			name: "plain resolver",
			schema: `
				schema {
					query: Query
				}

				type Query {
					hello: String
				}`,
			resolver: &traceErrPlainQuery{},
			query:    `{ hello }`,
			field:    "Query.hello",
			wantPath: []any{"hello"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rt := &fieldErrRecordingTracer{}
			opts := append([]graphql.SchemaOpt{graphql.Tracer(rt)}, tt.opts...)
			schema := graphql.MustParseSchema(tt.schema, tt.resolver, opts...)

			resp := schema.Exec(context.Background(), tt.query, "", nil)
			if len(resp.Errors) != 1 {
				t.Fatalf("got %d response errors, want 1: %v", len(resp.Errors), resp.Errors)
			}

			rt.mu.Lock()
			defer rt.mu.Unlock()

			got, called := rt.fieldErr[tt.field]
			if !called {
				t.Fatalf("TraceField finish func for %s was not called", tt.field)
			}
			if got == nil {
				t.Fatalf("TraceField finish func for %s was called with a nil error, want the resolver's error", tt.field)
			}
			if !reflect.DeepEqual(got.Path, tt.wantPath) {
				t.Errorf("error path = %v, want %v", got.Path, tt.wantPath)
			}
		})
	}
}
