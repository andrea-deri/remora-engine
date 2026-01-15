// Package btp defines the Benthos–Tidal Protocol contract.
package btp

import (
	"remora/pkg/protocol/transfer/btp"
)

// EmbeddedRegistry define the minimalist implementation of an embedded
// registry in certain persistence provider, used in monolithic applications
// or separated instance. The instance of [EmbeddedRegistry] directly uses
// the APIs implemented by provider implementation.

// EmbeddedRegistry provides a lightweight registry wrapper used when the
// persistence provider is embedded directly within the application (e.g. in
// monolithic deployments or distributed runtime instances). The registry does
// not introduce any additional transport layer: it forwards all operations
// directly to the underlying provider implementation.
//
// The embedded provider must implement the [btp.BenthosAPI] interface.
type EmbeddedRegistry struct {

	// BenthosAPI define the provider implementation to which all registry
	// operations are delegated.
	btp.BenthosAPI
}

// ReadAllReferences forwards the call to the underlying BenthosAPI implementation.
//
// For full details, refer to [BenthosAPI.ReadAllReferences].
func (registry *EmbeddedRegistry) ReadAllReferences() btp.ReferenceTransferData {

	return registry.BenthosAPI.ReadAllReferences()
}

// ReadSingleReference forwards the call to the underlying BenthosAPI implementation.
//
// For full details, refer to [BenthosAPI.ReadSingleReference].
func (registry *EmbeddedRegistry) ReadSingleReference(minnowId string) btp.ReferenceTransferData {

	return registry.BenthosAPI.ReadSingleReference(minnowId)
}

// ReadAllProfiles forwards the call to the underlying BenthosAPI implementation.
//
// For full details, refer to [BenthosAPI.ReadAllProfiles].
func (registry *EmbeddedRegistry) ReadAllProfiles() btp.ProfileTransferData {

	return registry.BenthosAPI.ReadAllProfiles()
}

// ReadSingleProfile forwards the call to the underlying BenthosAPI implementation.
//
// For full details, refer to [BenthosAPI.ReadSingleProfile].
func (registry *EmbeddedRegistry) ReadSingleProfile(minnowId string) btp.ProfileTransferData {

	return registry.BenthosAPI.ReadSingleProfile(minnowId)
}

// ReadAllEvents forwards the call to the underlying BenthosAPI implementation.
//
// For full details, refer to [BenthosAPI.ReadAllRReadAllEventseferences].
func (registry *EmbeddedRegistry) ReadAllEvents() btp.EventTransferData {

	return registry.BenthosAPI.ReadAllEvents()
}
