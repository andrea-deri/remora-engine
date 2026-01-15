package factory

import (
	"context"
	"fmt"

	btpimpl "remora/benthos/protocol/btp"
	"remora/benthos/registry"
	"remora/pkg/config"
	"remora/pkg/protocol/transfer/btp"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewRegistry creates a BenthosAPI registry instance for interacting with the persistence layer
// in order to retrieve Minnow Profiles. The function uses the provided configuration map to
// determine the type of storage backend and the communication mode, returning an implementation
// that satisfies the [btp.BenthosAPI] interface.
func NewRegistry(configMap config.ConfigMap) btp.BenthosAPI {

	// Determine the storage backend type
	var repository btp.BenthosAPI
	registryType := configMap.ReadString("REGISTRY_TYPE", "ND")
	switch registryType {

	case "filesystem":
		panic("An error occurred while opening connection to storage client. Type 'filesystem' currently not implemented")
	case "mongo":
		repository = newMongoRegistry(configMap)
	case "etcd":
		panic("An error occurred while opening connection to storage client. Type 'etcd' currently not implemented")
	default:
		panic("An error occurred while opening connection to storage client. Invalid registry type")
	}

	// Determine the registry API type (embedded vs external)
	var registry btp.BenthosAPI
	registryLinkType := configMap.ReadString("REGISTRY_API_TYPE", "embedded")
	switch registryLinkType {

	case "embedded":
		registry = &btpimpl.EmbeddedRegistry{BenthosAPI: repository}
	case "external":
		panic("An error occurred while opening connection to storage client. Invalid registry link type")
	}

	return registry
}

// newMongoRegistry creates a [btp.BenthosAPI] implementation backed by MongoDB.
// The function retrieves the connection string from the configuration map and establishes a
// connection to the MongoDB instance.
func newMongoRegistry(configMap config.ConfigMap) btp.BenthosAPI {

	storageConnectionString := configMap.ReadString("REGISTRY_CONNECTION_STRING", "-")
	if storageConnectionString == "-" {
		panic("An error occurred while opening connection to storage client. Parameter [REGISTRY_CONNECTION_STRING] required but not defined.")
	}

	log.Trace().Msgf("Trying to connect to MongoDB instance at [%v].", storageConnectionString)
	clientOptions := options.Client().ApplyURI(storageConnectionString)
	storageClient, error := mongo.Connect(context.TODO(), clientOptions)
	if error != nil {
		panic(fmt.Sprintf("An error occurred while opening connection to storage client [Connection string: %s]. %s", storageConnectionString, error))
	}

	// Generate the bundle attaching the storage client to the registry
	log.Trace().Msgf("Connected to MongoDB instance!")
	mongoRegistry := registry.MongoRegistry{
		StorageClient: storageClient,
	}
	return &mongoRegistry
}
