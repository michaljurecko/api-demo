package playerbiz

import "github.com/google/wire"

var WireSet = wire.NewSet(NewService)
