package configuration

import (
	"context"
	"sync"
)

// Store GCP Context for different environment
var gCPContext context.Context
var authOnce sync.Once

// GetContext will return Auth context of Google
func GetContext() context.Context {

	authOnce.Do(func() {
		gCPContext = context.Background()
	})
	return gCPContext
}
