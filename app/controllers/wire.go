package controllers

import (
	"github.com/google/wire"
)

var ProviderControllerSet = wire.NewSet(
	NewController,
)
