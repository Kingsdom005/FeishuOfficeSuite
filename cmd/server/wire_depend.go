//go:build wire
// +build wire

package data

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewData,
	NewDataWithMySQL,
	NewDataWithRedis,
	NewKafkaProducer,
	NewCasbinEnforcer,
)