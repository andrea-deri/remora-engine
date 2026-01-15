package proxy

import (
	"fmt"
	"net/http"
	"remora/harbor/protocol/message"
	"remora/pkg/config"
	"time"

	"github.com/rs/zerolog/log"
)

// HarborProxy manages incoming request made from client, redirecting them
// to [HarborDispatcher].
// It handles request content extraction and communication timeout.
type HarborProxy struct {

	// URL is the base path of the address on which the Harbor
	// Proxy service is esposed to external clients.
	BasePath string

	// DispatcherChannel is the channel that permits to send [message.HarborRequest],
	// generated from client HTTP requests, to Harbor Dispatcher in order to
	// start search and validation processes.
	// Currently, no size is defined for this channel, demanding the throttling
	// to the Harbor Dispatcher below.
	DispatcherChannel chan message.HarborRequest

	// StatusChannel defines the dedicated channel where the status change messages
	// are sent in order to being processed by [HarborProxy].
	StatusChannel chan string

	// Port is the value of the TCP port address on which the
	// Harbor Proxy service is exposed to external clients.
	Port uint

	// RequestTimeoutInSec is the time in seconds within which a request is
	// considered valid and after which it expires due to a timeout.
	RequestTimeoutInSec int
}

// NewProxy creates and returns a fully initialized HarborProxy instance.
// Configuration values are read from the provided ConfigMap, with defaults applied
// when a specific key is missing. The returned proxy is ready to be started.
func NewProxy(configMap config.ConfigMap, dispatcherChannel chan message.HarborRequest, statusChannel chan string) HarborProxy {

	return HarborProxy{
		Port:                uint(configMap.ReadInt("PROXY_EXPOSED_PORT", 8080)),
		BasePath:            configMap.ReadString("PROXY_PATH_PREFIX", "/"),
		DispatcherChannel:   dispatcherChannel,
		StatusChannel:       statusChannel,
		RequestTimeoutInSec: configMap.ReadInt("PROXY_REQUEST_TIMEOUT", 10),
	}
}

// StartListening starts the HTTP listener using the configured base path and exposed port.
// Incoming HTTP requests are processed by the first-level routing implemented by the proxy
// component, then forwarded through DispatcherChannel to the lower-level Harbor dispatcher
// for asynchronous analysis.
//
// This method is blocking and terminates only if http.ListenAndServe returns a fatal error.
// Until that moment, all requests are handled by the Proxy's RequestHandler instance.
func (proxy *HarborProxy) StartListening() {

	reqHandler := RequestHandler{
		BasePath:          proxy.getBasePath(),
		DispatcherChannel: proxy.DispatcherChannel,
		RequestTimeout:    time.Duration(proxy.RequestTimeoutInSec) * time.Second,
		StatusChannel:     proxy.StatusChannel,
	}
	reqHandler.Init()

	log.Info().
		Str("Component", "Harbor").
		Msgf("Application listening at path [%s] on port [%d]", proxy.getBasePath(), proxy.Port)

	// Register the handler and start the HTTP listener
	http.Handle(proxy.getBasePath(), &reqHandler)
	error := http.ListenAndServe(proxy.getPort(), nil)

	log.Info().
		Str("Component", "Harbor").
		Msgf("Application ended with status: [%s]", error)
}

// getBasePath returns the normalized base path for the Harbor Proxy.
// The returned value always starts and ends with a slash ('/'). If no base
// path is configured, the default root path ("/") is returned.
func (proxy *HarborProxy) getBasePath() string {

	basePath := proxy.BasePath
	if len(basePath) == 0 {
		return "/"
	}

	if basePath[0] != '/' {
		basePath = "/" + basePath
	}

	if basePath[len(basePath)-1] != '/' {
		basePath = basePath + "/"
	}
	return basePath
}

// getPort returns the formatted port string used by the Harbor Proxy listener.
// The returned value is always prefixed with ':' as required by http.ListenAndServe.
// If no port is configured (default int value), the default port 8080 is used.
func (proxy *HarborProxy) getPort() string {

	port := proxy.Port
	if port == 0 {
		port = 8080
	}
	return ":" + fmt.Sprint(port)
}
