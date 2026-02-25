package cli

import (
	"context"

	"github.com/jackmorganxyz/projectsCLI/internal/config"
)

type runtimeContextKey struct{}

// RuntimeContext carries root-level runtime state to all subcommands.
type RuntimeContext struct {
	Config     config.Config
	ConfigPath string
	JSON       bool
}

// WithRuntimeContext attaches runtime state to a context.
func WithRuntimeContext(ctx context.Context, runtime RuntimeContext) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, runtimeContextKey{}, runtime)
}

// RuntimeFromContext extracts runtime state from context.
func RuntimeFromContext(ctx context.Context) (RuntimeContext, bool) {
	if ctx == nil {
		return RuntimeContext{}, false
	}
	runtime, ok := ctx.Value(runtimeContextKey{}).(RuntimeContext)
	return runtime, ok
}
