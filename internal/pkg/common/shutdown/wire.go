package shutdown

import "github.com/google/wire"

var WireSet = wire.NewSet(NewStack)
