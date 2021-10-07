package middleware

import (
	"context"
	"fmt"
	"log"
	"sync"
)

var middlewareManager *MiddlewareManager

func init() {
	once := sync.Once{}
	once.Do(func() {
		GetMiddlewareManager()
	})
}

// GetMiddlewareManager return a global middleware manager struct (singleton).
func GetMiddlewareManager() *MiddlewareManager {
	if middlewareManager == nil {
		middlewareManager = &MiddlewareManager{middlewareMap: sync.Map{}}
	}
	return middlewareManager
}

// Register the middleware plugin in the manager.
func (m *MiddlewareManager) Register(name string, action MiddlewareHandler) {
	log.Println("Adding middleware", name, action)
	m.middlewareMap.Store(name, action)
}

// HasMiddleware Check if the manager has the middleware stored.
func (m *MiddlewareManager) HasMiddleware(name string) bool {
	_, isFound := m.middlewareMap.Load(name)
	if !isFound {
		return false
	}
	return true
}

// GetMiddleware gets a middleware stored in the map.
func (m *MiddlewareManager) GetMiddleware(name string) (MiddlewareHandler, error) {
	middleware, isFound := m.middlewareMap.Load(name)
	if !isFound {
		return nil, fmt.Errorf("middleware %s not found", name)
	}
	return middleware.(MiddlewareHandler), nil
}

// GetAllMiddleware gets all the middleware stored in the manager.
func (m *MiddlewareManager) GetAllMiddleware() []MiddlewareHandler {
	var middlewares []MiddlewareHandler
	m.middlewareMap.Range(func(key, value interface{}) bool {
		middlewares = append(middlewares, value.(MiddlewareHandler))
		return true
	})
	return middlewares
}

// Chain the middleware in consecutive order. This is useful for processing requests depending on the business constraints.
// Returns an error if it occurred during execution.
func (m *MiddlewareManager) Chain(ctx context.Context, middleware ...string) error {
	var finalMiddleware MiddlewareHandler
	if len(middleware) == 0 {
		return nil
	}

	for i, key := range middleware {
		if !m.HasMiddleware(key) {
			return fmt.Errorf("middleware %s not found", key)
		}

		m, _ := m.middlewareMap.Load(key)
		if i == 0 {
			// assign the first instance
			finalMiddleware = m.(MiddlewareHandler)
		} else {
			// try to chain the middleware, stop execution if an error occurs
			mw, err := finalMiddleware.Chain(ctx, m.(MiddlewareHandler))
			if err != nil {
				return err
			}
			finalMiddleware = mw
		}
	}

	return finalMiddleware.Execute(ctx)
}
