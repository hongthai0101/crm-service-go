package controller_inside

import "github.com/google/wire"

var ProviderInsideControllerSet = wire.NewSet(
	NewInsideController,
)
