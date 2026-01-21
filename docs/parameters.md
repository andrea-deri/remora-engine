---
title: Parameters and configuration
nav_order: 3
has_children: true
---

This section lists the configuration parameters required by the REMORA Engine to initialize and operate correctly. 

### 3.1. Proxy and Dispatcher

| Key | Description | Accepted values | Default value |
|-----|--------------|----------------|----------------|
| `PROXY_EXPOSED_PORT` | Port on which the proxy listens for incoming requests. |  | `8080` |
| `PROXY_PATH_PREFIX` | URL path prefix handled by the proxy. |  | `/` |
| `PROXY_REQUEST_TIMEOUT` | Maximum time (in seconds) to wait before timing out a request. |  | `10` |
| `DISPATCHER_MAX_ACCEPTED_REQUESTS` | Maximum number of accepted requests before returning an HTTP “Too Many Requests” error. |  | `64` |
| `DISPATCHER_MAX_CONCURRENT_REQUESTS` | Maximum number of concurrent requests processed in parallel.<br><u>Currently not implemented.</u> |  | `64` |
| `DISPATCHER_CAN_CACHE_CONTENT` | Enabling whether responses can be cached.<br><u>Currently not implemented.</u> | `true`<br>`false` | `false` |
| `DISPATCHER_SPECIAL_ENDPOINT_PREFIX` | URL path prefix after `PROXY_PATH_PREFIX` used for routing special endpoints |  | `r` |

### 3.2. Registries and Persistence

| Key | Description | Accepted values | Default value |
|-----|--------------|----------------|----------------|
| `REGISTRY_TYPE` | Type of persistence used for registry. | `filesystem`<br>`mongo`<br>`etcd` | `ND` |
| `REGISTRY_API_TYPE` | Type of link between Tidal engine and Benthos registry. | `embedded`<br>`external` | `embedded` |
| `REGISTRY_CONNECTION_STRING` | Connection string for the persistence layer. | For MongoDB: `mongodb://<host>`<br>For filesystem: `./path/to/folder` | `-` |
| `TIDAL_MINNOW_PROFILES_UPDATE_POLLING_SECONDS` | Polling interval in seconds between attempts to update Minnow profiles in Engine. |  | `60` |
| `REGISTRY_LOAD_ALL_PROFILES_AT_STARTUP` | Enable preloading for profiles. If `true`, execute the loading of all Minnow profiles at startup (faster runtime). If `false`, loads them lazily (faster startup).<br><u>Currently preload is default implementation.</u> | `true`<br>`false` | `true` |


### 3.3. Logging Strategies

| Key | Description | Accepted values | Default value |
|-----|--------------|----------------|----------------|
| `HARBOR_DISPATCHER_LOGGING_LEVEL` | Logging verbosity for the dispatcher. | `trace`<br>`debug`<br>`info`<br>`warn`<br>`error`<br>`fatal`<br>`panic`<br>`disabled` | `info` |
| `HARBOR_PROXY_LOGGING_LEVEL` | Logging verbosity for the proxy. | `trace`<br>`debug`<br>`info`<br>`warn`<br>`error`<br>`fatal`<br>`panic`<br>`disabled` | `info` |
