// Package btp defines the Benthos–Tidal Protocol contract.
package btp

// BenthosAPI defines the guidelines on which all the Benthos implementations
// must follows in order to exposes the required behaviors on persistence level.
type BenthosAPI interface {

	// ReadAllReferences retrieves all Minnow references from the persistence layer.
	ReadAllReferences() ReferenceTransferData

	// ReadSingleReference retrieves a single Minnow reference from the persistence layer.
	ReadSingleReference(minnowId string) ReferenceTransferData

	// ReadAllProfiles retrieves all Minnow profiles from the persistence layer.
	ReadAllProfiles() ProfileTransferData

	// ReadSingleProfile retrieves a single Minnow profile from the persistence layer.
	ReadSingleProfile(minnowId string) ProfileTransferData

	// ReadAllEvents retrieves all events on Minnow profiles from the persistence layer.
	ReadAllEvents() EventTransferData
}
