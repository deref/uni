package internal

import "github.com/evanw/esbuild/pkg/api"

// WARNING: Temporarily turns these loaders in to effective no-ops.
// TODO: Plugins or something, but probably not that.
var loaders = map[string]api.Loader{
	".scss": api.LoaderText,
	".svg":  api.LoaderText,
}
