package fw

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
	Path    string
	Page    []byte
}

type GraphGopherConfig struct {
	GraphqlPath string
	G           GraphQLAPI
	Graphiql    GraphiQlConfig
}
