package service_inside

import (
	"github.com/google/wire"
)

var ProviderInsideServiceSet = wire.NewSet(
	NewInsideService,
)
