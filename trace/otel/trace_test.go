package otel_test

import (
	"testing"

	"go.opentelemetry.io/otel"

	"github.com/tribunadigital/graphql-go"
	"github.com/tribunadigital/graphql-go/example/starwars"
	otelgraphql "github.com/tribunadigital/graphql-go/trace/otel"
	"github.com/tribunadigital/graphql-go/trace/tracer"
)

func TestInterfaceImplementation(t *testing.T) {
	var _ tracer.ValidationTracer = &otelgraphql.Tracer{}
	var _ tracer.Tracer = &otelgraphql.Tracer{}
}

func TestTracerOption(t *testing.T) {
	_, err := graphql.ParseSchema(starwars.Schema, nil, graphql.Tracer(otelgraphql.DefaultTracer()))
	if err != nil {
		t.Fatal(err)
	}

	_, err = graphql.ParseSchema(starwars.Schema, nil, graphql.Tracer(&otelgraphql.Tracer{Tracer: otel.Tracer("example")}))
	if err != nil {
		t.Fatal(err)
	}
}
