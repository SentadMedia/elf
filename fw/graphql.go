package fw

import "net/http"

type GraphQlPath string
type RelayHandler http.Handler

type Resolver interface{}

type GraphQLAPI interface {
	GetSchema() string
	GetResolver() Resolver
}

type GraphQLScalar interface {
	ImplementsGraphQLType(name string) bool
	UnmarshalGraphQL(input interface{}) error
	MarshalJSON() ([]byte, error)
}

type GraphiQlConfig struct {
	Include bool
	Path    GraphQlPath
	Page    []byte
}

type GraphGopherConfig struct {
	GraphqlPath string
	G           GraphQLAPI
	Graphiql    GraphiQlConfig
}
