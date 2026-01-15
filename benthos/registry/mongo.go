package registry

import (
	"context"
	"remora/benthos/entity"
	"remora/pkg/customerror"
	"remora/pkg/protocol/transfer/btp"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoRegistry represents the concrete implementation of registry's processing
// APIs backed by MongoDB for data persistence.
type MongoRegistry struct {

	// StorageClient defines the MongoDB client responsible for retrieving data
	// from the document-oriented database.
	StorageClient *mongo.Client
}

// ReadAllReferences retrieves all Minnow references from the MongoDB persistence layer.
// It queries the `minnows` collection in the `remora` database and returns the results as a slice
// of [btp.ReferenceTransferData]. The operation is executed with a 10-second timeout context.
// If any error occurs during lookup or decoding, the method returns a response with error code.
func (registry *MongoRegistry) ReadAllReferences() btp.ReferenceTransferData {

	var response btp.ReferenceTransferData
	profileEntities, errorCode := registry.readAllMinnows()
	if errorCode != "" {

		response.Error = string(errorCode)
		return response
	}

	// Converting all found entities as reference
	references := make([]btp.MinnowReference, 0)
	for _, profileEntity := range profileEntities {

		extractedReference := btp.MinnowReference{
			Id:             profileEntity.Id,
			Method:         profileEntity.Method,
			Path:           profileEntity.Path,
			SpecialHeaders: profileEntity.SpecialHeaders,
		}
		references = append(references, extractedReference)
	}
	response.References = references
	response.Size = len(response.References)
	return response
}

// ReadAllProfiles retrieves all Minnow profiles from the MongoDB persistence layer.
// It queries the `minnows` collection in the `remora` database and returns the results as a slice
// of [btp.ProfileTransferData]. The operation is executed with a 10-second timeout context.
// If any error occurs during query or decoding, the method returns a response with error cause.
func (registry *MongoRegistry) ReadAllProfiles() btp.ProfileTransferData {

	var response btp.ProfileTransferData
	profileEntities, errorCode := registry.readAllMinnows()
	if errorCode != "" {

		response.Error = string(errorCode)
		return response
	}

	// Converting all found entities as Profiles
	profiles := make([]btp.MinnowProfile, 0)
	for _, profileEntity := range profileEntities {

		var behaviors []btp.MinnowBehavior
		for _, behaviorEntity := range profileEntity.Behaviors {

			behaviorEffect := behaviorEntity.Effect
			behavior := btp.MinnowBehavior{
				Condition: behaviorEntity.Condition,
				Effect: btp.MinnowBehaviorEffect{
					Type:         behaviorEffect.Type,
					StatusCode:   behaviorEffect.StatusCode,
					Headers:      behaviorEffect.Headers,
					ContentType:  behaviorEffect.ContentType,
					RawBody:      behaviorEffect.RawBody,
					MappableBody: behaviorEffect.MappableBody,
					Execute:      behaviorEffect.Execute,
				},
			}
			behaviors = append(behaviors, behavior)
		}

		extractedProfile := btp.MinnowProfile{
			Id:             profileEntity.Id,
			Method:         profileEntity.Method,
			Path:           profileEntity.Path,
			SpecialHeaders: profileEntity.SpecialHeaders,
			Behaviors:      behaviors,
		}
		profiles = append(profiles, extractedProfile)
	}
	response.Profiles = profiles
	response.Size = len(response.Profiles)
	return response
}

// ReadAllEvents retrieves all Minnow events from the MongoDB persistence layer.
// It queries the `minnow_events` collection in the `remora` database and returns the results as a slice
// of [btp.EventTransferData]. The operation is executed with a 10-second timeout context.
// If any error occurs during query or decoding, the method returns a response with error cause.
func (registry *MongoRegistry) ReadAllEvents() btp.EventTransferData {

	var response btp.EventTransferData
	eventEntities, errorCode := registry.readAllEvents()
	if errorCode != "" {

		response.Error = string(errorCode)
		return response
	}

	var eventsToDelete []primitive.ObjectID

	events := make([]btp.MinnowEvent, 0)
	for _, eventEntity := range eventEntities {

		// Include the event ID as MongoDB's native ObjectID in order to
		// delete that document after being consumed
		eventObjectId, err := primitive.ObjectIDFromHex(eventEntity.EventId)
		if err != nil {
			log.Error().Msgf("An error occurred while trying to generate Mongo ObjectID from event id [%s]. %s", eventEntity.EventId, err)
		} else {
			eventsToDelete = append(eventsToDelete, eventObjectId)
		}

		// The consumed event is converted in structure that will be returned
		extractedEvent := btp.MinnowEvent{
			MinnowId: eventEntity.MinnowId,
			Operation: btp.MinnowEventOperation{
				Type: eventEntity.Operation.Type,
				Reference: btp.MinnowEventReference{
					Method:         eventEntity.Operation.Reference.Method,
					Path:           eventEntity.Operation.Reference.Path,
					SpecialHeaders: eventEntity.Operation.Reference.SpecialHeaders,
				},
			},
		}
		events = append(events, extractedEvent)
	}
	response.Events = events

	// Delete all computed events in order to avoid re-applying them
	if len(eventsToDelete) > 0 {

		registry.deleteEvents(eventsToDelete)
	}
	return response
}

// ReadSingleReference retrieves a single Minnow references from the MongoDB persistence layer,
// searching by its identifier. It queries the `minnows` collection in the `remora` database and
// returns the results as a single [btp.ReferenceTransferData]. The operation is executed with
// a 10-second timeout context. If any error occurs during query or decoding, the method returns
// a response with error cause.
func (registry *MongoRegistry) ReadSingleReference(minnowId string) btp.ReferenceTransferData {

	var response btp.ReferenceTransferData
	profileEntity, errorCode := registry.readSingleMinnow(minnowId)
	if errorCode != "" {

		response.Error = string(errorCode)
		return response
	}

	reference := btp.MinnowReference{
		Id:             profileEntity.Id,
		Method:         profileEntity.Method,
		Path:           profileEntity.Path,
		SpecialHeaders: profileEntity.SpecialHeaders,
	}
	response.References = append(response.References, reference)
	response.Size = len(response.References)

	return response
}

// ReadSingleProfile retrieves a single Minnow profile from the MongoDB persistence layer,
// searching by its identifier. It queries the `minnows` collection in the `remora` database and
// returns the results as a single [btp.ProfileTransferData]. The operation is executed with
// a 10-second timeout context. If any error occurs during query or decoding, the method returns
// a response with error cause.
func (registry *MongoRegistry) ReadSingleProfile(minnowId string) btp.ProfileTransferData {

	var response btp.ProfileTransferData
	profileEntity, errorCode := registry.readSingleMinnow(minnowId)
	if errorCode != "" {

		response.Error = string(errorCode)
		return response
	}

	// Extract all behaviors from original entity for final response
	var behaviors []btp.MinnowBehavior
	for _, behaviorEntity := range profileEntity.Behaviors {

		behaviorEntityEffect := behaviorEntity.Effect
		behavior := btp.MinnowBehavior{
			Condition: behaviorEntity.Condition,
			Effect: btp.MinnowBehaviorEffect{
				Type:         behaviorEntityEffect.Type,
				StatusCode:   behaviorEntityEffect.StatusCode,
				Headers:      behaviorEntityEffect.Headers,
				ContentType:  behaviorEntityEffect.ContentType,
				RawBody:      behaviorEntityEffect.RawBody,
				MappableBody: behaviorEntityEffect.MappableBody,
				Execute:      behaviorEntityEffect.Execute,
			},
		}
		behaviors = append(behaviors, behavior)
	}

	profile := btp.MinnowProfile{
		Id:             profileEntity.Id,
		Method:         profileEntity.Method,
		Path:           profileEntity.Path,
		SpecialHeaders: profileEntity.SpecialHeaders,
		Behaviors:      behaviors,
	}
	response.Profiles = append(response.Profiles, profile)
	response.Size = len(response.Profiles)

	return response
}

// readAllMinnows retrieves every MinnowProfile document from the MongoDB collection.
//
// It returns either the full result set or a specific error code describing the failure cause.
func (registry *MongoRegistry) readAllMinnows() ([]entity.MinnowProfile, customerror.RemoraErrorCode) {

	log.Trace().Msgf("Reading all Minnow entities...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Empty filter ensures that all documents are fetched.
	filter := bson.D{}

	collection := registry.StorageClient.Database("remora").Collection("minnows")
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Error().Msgf("An error occurred while trying to find Minnow entities from DB. %s", err)
		return nil, customerror.CODE_BENTHOS_MONGODB_PROFILES_FIND_FAILED
	}

	// Ensures cursor resources are released in all execution paths.
	defer cursor.Close(ctx)

	var profileEntities []entity.MinnowProfile
	err = cursor.All(ctx, &profileEntities)
	if err != nil {
		log.Error().Msgf("An error occurred while trying to scan Minnow entities from cursor. %s", err)
		return nil, customerror.CODE_BENTHOS_MONGODB_PROFILES_FIND_FAILED
	}
	log.Trace().Msgf("Minnow entities read completed! [%v]", profileEntities)

	return profileEntities, ""
}

// readAllEvents retrieves every MinnowEvent document from the collection.
//
// It returns either the full result set or an error code describing the failure cause.
func (registry *MongoRegistry) readAllEvents() ([]entity.MinnowEvent, customerror.RemoraErrorCode) {

	log.Trace().Msgf("Reading all Minnow events...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Empty filter ensures that all documents are fetched.
	filter := bson.D{}

	collection := registry.StorageClient.Database("remora").Collection("minnow_events")
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Error().Msgf("An error occurred while trying to find Minnow events from DB. %s", err)
		return nil, customerror.CODE_BENTHOS_MONGODB_EVENTS_FIND_FAILED
	}

	// Ensures cursor resources are released in all execution paths.
	defer cursor.Close(ctx)

	var eventEntities []entity.MinnowEvent
	err = cursor.All(ctx, &eventEntities)
	if err != nil {
		log.Error().Msgf("An error occurred while trying to retrieve all Minnow events by cursor from DB. %s", err)
		return nil, customerror.CODE_BENTHOS_MONGODB_EVENTS_FIND_FAILED
	}
	log.Trace().Msgf("Minnow events read completed! [%v]", eventEntities)

	return eventEntities, ""
}

// readSingleMinnow retrieves a MinnowProfile by its unique ID.
//
// It returns either the decoded entity or a specific error code indicating the failure reason.
func (registry *MongoRegistry) readSingleMinnow(minnowId string) (entity.MinnowProfile, customerror.RemoraErrorCode) {

	log.Trace().Msgf("Reading single Minnow entity with ID [%s]...", minnowId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := registry.StorageClient.Database("remora").Collection("minnows")

	var profileEntity entity.MinnowProfile
	filter := bson.D{{Key: "id", Value: minnowId}}
	err := collection.FindOne(ctx, filter).Decode(&profileEntity)

	switch {

	case err == mongo.ErrNoDocuments:
		log.Error().Msgf("An error occurred while trying to find Minnow entity from DB. No entity found with ID [%s]", minnowId)
		return entity.MinnowProfile{}, customerror.CODE_BENTHOS_ENTITY_NOT_FOUND

	case err != nil:
		log.Error().Msgf("An error occurred while trying to retrieve Minnow entity from DB. %s", err)
		return entity.MinnowProfile{}, customerror.CODE_BENTHOS_MONGODB_PROFILES_FIND_FAILED
	}

	log.Trace().Msgf("Minnow entity with ID [%s] found! [%v]", minnowId, profileEntity)
	return profileEntity, ""
}

// deleteEvents removes all MinnowEvent documents whose IDs match the provided list.
//
// It returns an error code only when the deletion operation fails.
func (registry *MongoRegistry) deleteEvents(events []primitive.ObjectID) customerror.RemoraErrorCode {

	log.Trace().Msgf("Deleting [%v] consumed Minnow events...", len(events))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := registry.StorageClient.Database("remora").Collection("minnow_events")

	deleteFilter := bson.M{"_id": bson.M{"$in": events}}
	result, err := collection.DeleteMany(ctx, deleteFilter)

	if err != nil {
		log.Error().Msgf("An error occurred while trying to delete Minnow events from DB. %s", err)
		return customerror.CODE_BENTHOS_MONGODB_EVENTS_DELETE_FAILED
	}

	log.Trace().Msgf("Required delete of [%v] Minnow events. Performed delete of [%v] events. %v", len(events), result.DeletedCount, err)
	return ""
}
