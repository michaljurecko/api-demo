package distlock

import "github.com/google/wire"

var WireSet = wire.NewSet(NewLocker)
