package main

import (
	"fmt"

	benthosfactory "remora/benthos/factory"
	"remora/benthos/protocol/btp"
	"remora/harbor/dispatcher"
	"remora/harbor/proxy"
	"remora/pkg/config"
	tidalfactory "remora/tidal/factory"
	"remora/tidal/protocol/htp"

	"github.com/rs/zerolog/log"
)

// VERSION holds the current version of the REMORA application.
var VERSION = "undefined"

// CONFIG_FILE_PATH defines the path to the file that includes
// Engine's environment variables
var CONFIG_FILE_PATH = "properties.env"

// main is the entry point for the REMORA application monolith.
// It loads configuration, connects to MongoDB, initializes engines and
// starts the dispatcher and proxy.
func main() {

	log.Info().Msgf("##### REMORA - Version [%s] #####", VERSION)

	configMap, error := config.LoadConfig(CONFIG_FILE_PATH)
	if error != nil {
		panic(fmt.Sprintf("An error occurred while loading configuration properties: [%s]", error))
	}

	applicationStatusChannel := make(chan string)

	registry := benthosfactory.NewRegistry(configMap)
	embeddedRegistry := registry.(*btp.EmbeddedRegistry)

	engine := tidalfactory.NewEmbeddedEngine(configMap, registry)
	embeddedEngine := engine.(*htp.EmbeddedEngine)
	embeddedEngine.StatusChannel = applicationStatusChannel
	go embeddedEngine.Init(configMap, embeddedRegistry)

	// Start the HarborDispatcher and manage new incoming requests
	dispatcher := dispatcher.NewDispatcher(configMap, embeddedEngine)
	dispatcher.Start()

	// Start the HarborProxy and start accepting incoming requests
	proxy := proxy.NewProxy(configMap, dispatcher.InputChannel, applicationStatusChannel)
	proxy.StartListening()
}
