package middleware

import (
	"context"
	"sync"
)

type (
	// MiddlewareHandler is the interface that needs to be implemented in order to make functioning middleware
	MiddlewareHandler interface {

		// Execute should be called in the last
		Execute(ctx context.Context) error

		// Chain does some logic and returns the next middleware in the chain and an error, if one occurs during execution
		Chain(ctx context.Context, next MiddlewareHandler) (MiddlewareHandler, error)
	}

	// MiddlewareManager registers, stores and chains middleware.
	MiddlewareManager struct {
		middlewareMap sync.Map
	}
)
