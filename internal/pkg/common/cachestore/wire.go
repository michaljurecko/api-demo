package cachestore

import "github.com/google/wire"

var WireSet = wire.NewSet(New)
