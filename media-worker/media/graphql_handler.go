package media

import (
	"context"
	"net/http"

	"github.com/machinebox/graphql"
)

type headerCapturingTransport struct {
	underlyingTransport http.RoundTripper
	headers             http.Header
}

func (h *headerCapturingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := h.underlyingTransport.RoundTrip(req)
	if err == nil {
		h.headers = resp.Header
	}
	return resp, err
}

type GraphqlHandler struct {
	client    *graphql.Client
	transport *headerCapturingTransport
}

func NewGraphQLHandler(url string) *GraphqlHandler {
	capturingTransport := &headerCapturingTransport{
		underlyingTransport: http.DefaultTransport,
	}

	httpClient := &http.Client{Transport: capturingTransport}
	graphqlClient := graphql.NewClient(url, graphql.WithHTTPClient(httpClient))

	return &GraphqlHandler{
		client:    graphqlClient,
		transport: capturingTransport,
	}
}

func (handler *GraphqlHandler) Query(query string, variables map[string]interface{}, graphqlResponse interface{}) (headers http.Header, err error) {
	graphqlRequest := graphql.NewRequest(query)

	for key, value := range variables {
		graphqlRequest.Var(key, value)
	}

	if err := handler.client.Run(context.Background(), graphqlRequest, graphqlResponse); err != nil {
		return handler.transport.headers, err
	}

	return handler.transport.headers, nil
}
