---
title: Parameters and configuration
nav_order: 3
has_children: true
---

This section lists the configuration parameters required by the REMORA Engine to initialize and operate correctly. 

### 3.1. Proxy and Dispatcher <a name="configuration-parameters-proxy"></a>

| Key | Description | Default value |
|-----|--------------|----------------|
| `PROXY_EXPOSED_PORT` | Port on which the proxy listens for incoming requests. | `8080` |
| `PROXY_PATH_PREFIX` | URL path prefix handled by the proxy. | `/` |
| `PROXY_REQUEST_TIMEOUT` | Maximum time (in seconds) to wait before timing out a request. | `10` |
| `DISPATCHER_MAX_ACCEPTED_REQUESTS` | Maximum number of accepted requests before returning an HTTP “Too Many Requests” error. | `64` |
| `DISPATCHER_MAX_CONCURRENT_REQUESTS` | Maximum number of concurrent requests processed in parallel.<br><u>Currently not implemented.</u> | `64` |
| `DISPATCHER_CAN_CACHE_CONTENT` | Enabling whether responses can be cached.<br><u>Currently not implemented.</u> | `false` |
| `DISPATCHER_SPECIAL_ENDPOINT_PREFIX` | URL path prefix after `PROXY_PATH_PREFIX` used for routing special endpoints | `r` |

### 3.2. Registries and Persistence <a name="configuration-parameters-registry"></a>

| Key | Description | Default value |
|-----|--------------|----------------|
| `REGISTRY_TYPE` | Type of persistence used for registry. Supported types: <ul><li>`filesystem`</li><li>`mongo`</li><li>`etcd`</li></ul> | `ND` |
| `REGISTRY_API_TYPE` | Type of link between Tidal engine and Benthos registry. Supported types: <ul><li>`embedded`</li><li>`external`</li></ul> | `embedded` |
| `REGISTRY_CONNECTION_STRING` | Connection string for the persistence layer. Standard strings:<ul><li>MongoDB: `mongodb://<host>`</li><li>Filesystem: `./path/to/folder`</li></ul> | `-` |
| `TIDAL_MINNOW_PROFILES_UPDATE_POLLING_SECONDS` | Polling interval in seconds between attempts to update Minnow profiles in Engine. | `60` |
| `REGISTRY_LOAD_ALL_PROFILES_AT_STARTUP` | Enable preloading for profiles. If `true`, execute the loading of all Minnow profiles at startup (faster runtime). If `false`, loads them lazily (faster startup).<br><u>Currently preload is default implementation.</u> | `true` |


### 3.3. Logging Strategies <a name="configuration-parameters-logs"></a>

| Key | Description | Default value |
|-----|--------------|----------------|
| `HARBOR_DISPATCHER_LOGGING_LEVEL` | Logging verbosity for the dispatcher. Supported levels:<ul><li>`trace`</li> <li>`debug`</li><li>`info`</li><li>`warn`</li><li>`error`</li><li>`fatal`</li><li>`panic`</li><li>`disabled`</li></ul> | `info` |
| `HARBOR_PROXY_LOGGING_LEVEL` | Logging verbosity for the proxy. Supported levels:<ul><li>`trace`</li> <li>`debug`</li><li>`info`</li><li>`warn`</li><li>`error`</li><li>`fatal`</li><li>`panic`</li><li>`disabled`</li></ul> | `info` |
