package factory

import (
	"remora/pkg/config"
	"remora/pkg/protocol/transfer/btp"
	"remora/tidal/engine"
	"remora/tidal/protocol/htp"
)

// NewEmbeddedEngine creates and returns a new TidalEngine instance configured as an embedded component.
// The returned engine is ready to be initialized and started. The function accepts a configuration map
// and a Benthos API interface, which are used internally by the engine for initialization and data access.
func NewEmbeddedEngine(configMap config.ConfigMap, benthosAPI btp.BenthosAPI) engine.Engine {

	return htp.NewEmbeddedEngine(&engine.TidalEngine{})
}
